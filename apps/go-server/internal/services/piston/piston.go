// Package piston ports services/piston.py: a pooled HTTP client that submits
// code to a Piston-compatible execution server and normalizes its response.
package piston

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// langInfo is (piston language name, version, file extension).
type langInfo struct {
	name, version, ext string
}

var languageMap = map[string]langInfo{
	"PYTHON":     {"python", "3.10.0", "py"},
	"JAVASCRIPT": {"javascript", "18.15.0", "js"},
	"JAVA":       {"java", "15.0.2", "java"},
	"C":          {"c", "10.2.0", "c"},
	"CPP":        {"c++", "10.2.0", "cpp"},
}

var languageAliases = map[string]string{
	"PY": "PYTHON", "JS": "JAVASCRIPT", "NODE": "JAVASCRIPT", "C++": "CPP", "CXX": "CPP",
}

// NormalizeLanguage ports the alias-resolution used by both this client and
// the assessment router's template lookup.
func NormalizeLanguage(language string) string {
	normalized := strings.ToUpper(strings.TrimSpace(language))
	if alias, ok := languageAliases[normalized]; ok {
		return alias
	}
	return normalized
}

// Client talks to a Piston-compatible execution server. One shared instance
// (single pooled *http.Client) should be reused across requests rather than
// constructing a client per call.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New builds a Client for the given CODE_EXECUTION_SERVER base URL.
func New(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Result mirrors the dict returned by execute_code: either the raw Piston
// response (run/compile/language/version) or {"error", "details"}.
type Result map[string]any

// Stdout extracts run.stdout with the same fallbacks as _extract_stdout.
func (r Result) Stdout() string {
	if run, ok := r["run"].(map[string]any); ok {
		if s, ok := run["stdout"].(string); ok && s != "" {
			return s
		}
	}
	if s, ok := r["stdout"].(string); ok && s != "" {
		return s
	}
	if s, ok := r["output"].(string); ok {
		return s
	}
	return ""
}

// Error returns the error message, if the result represents a failure.
func (r Result) Error() string {
	if e, ok := r["error"].(string); ok {
		return e
	}
	return ""
}

func (r Result) Details() string {
	if d, ok := r["details"].(string); ok {
		return d
	}
	return ""
}

func candidateURLs(baseURL string) []string {
	switch {
	case strings.HasSuffix(baseURL, "/api/v2/execute"):
		return []string{baseURL}
	case strings.HasSuffix(baseURL, "/api/v2/piston"):
		return []string{baseURL + "/execute"}
	case strings.HasSuffix(baseURL, "/api/v2"):
		return []string{baseURL + "/execute", baseURL + "/piston/execute"}
	default:
		return []string{baseURL + "/api/v2/piston/execute", baseURL + "/api/v2/execute"}
	}
}

// Execute submits code for the given language ("PYTHON", "JAVASCRIPT", ...;
// aliases like "JS"/"C++" are normalized) and returns the Piston result.
func (c *Client) Execute(ctx context.Context, language, code, stdin string) Result {
	lang := NormalizeLanguage(language)
	info, ok := languageMap[lang]
	if !ok {
		return Result{"error": fmt.Sprintf("Unsupported language: %s", language)}
	}

	payload := map[string]any{
		"language": info.name, "version": info.version,
		"files":                []map[string]string{{"name": "main." + info.ext, "content": code}},
		"stdin":                stdin,
		"args":                 []string{},
		"compile_timeout":      3000,
		"run_timeout":          3000,
		"compile_memory_limit": -1,
		"run_memory_limit":     -1,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return Result{"error": fmt.Sprintf("failed to encode request: %v", err)}
	}

	var lastErr Result
	for _, target := range candidateURLs(c.baseURL) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(body))
		if err != nil {
			lastErr = Result{"error": fmt.Sprintf("Execution server request failed: %v", err)}
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = c.connectionError(target, err)
			continue
		}

		result, shouldReturn := handleResponse(resp)
		if shouldReturn {
			return result
		}
		lastErr = result
	}

	if lastErr != nil {
		return lastErr
	}
	return Result{"error": "Execution server request failed", "details": "No reachable execution endpoint found"}
}

func (c *Client) connectionError(target string, err error) Result {
	u, parseErr := url.Parse(target)
	if parseErr == nil {
		host := u.Hostname()
		isLocal := host == "localhost" || host == "127.0.0.1" || host == "piston"
		port := u.Port()
		if isLocal && (port == "" || port == "2000") {
			return Result{"error": fmt.Sprintf(
				"Execution server request failed: local Piston is not reachable at %s. "+
					"Start it with `docker compose up -d piston` or set CODE_EXECUTION_SERVER "+
					"to a reachable Piston instance.", target,
			)}
		}
	}
	return Result{"error": fmt.Sprintf("Execution server request failed: %v", err)}
}

// handleResponse reads resp and returns (result, shouldReturn).
// shouldReturn=false only for 404/405, meaning the caller should try the
// next candidate URL rather than surfacing this as the final result.
func handleResponse(resp *http.Response) (Result, bool) {
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		details := readLimited(resp.Body, 500)
		message := "Execution server is not authorized for this project."
		if strings.Contains(strings.ToLower(details), "whitelist") {
			message = "Execution server access denied (EMKC now requires whitelist). " +
				"Set CODE_EXECUTION_SERVER to your own Piston instance."
		}
		return Result{"error": message, "details": details}, true
	}

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
		details := readLimited(resp.Body, 500)
		return Result{
			"error":   fmt.Sprintf("Execution server returned HTTP %d", resp.StatusCode),
			"details": details,
		}, false
	}

	if resp.StatusCode >= 400 {
		details := readLimited(resp.Body, 500)
		return Result{
			"error":   fmt.Sprintf("Execution server returned HTTP %d", resp.StatusCode),
			"details": details,
		}, true
	}

	var result Result
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		return Result{
			"error":   "Execution server returned invalid response format",
			"details": readLimited(resp.Body, 500),
		}, true
	}
	return result, true
}

func readLimited(r io.Reader, n int64) string {
	b, _ := io.ReadAll(io.LimitReader(r, n))
	return string(b)
}
