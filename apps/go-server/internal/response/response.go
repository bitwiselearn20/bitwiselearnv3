// Package response reproduces the exact JSON envelope of the legacy
// Python utils/api_response.py so the Next.js frontend needs zero changes.
//
// Envelope shape: {"statusCode", "message", "data", "error"}.
// Keys inside `data` are recursively converted from snake_case to camelCase,
// matching the Python _camel_keys behaviour.
package response

import (
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

// Envelope is the standard API response body.
type Envelope struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
	Error      any    `json:"error"`
}

// toCamel converts snake_case to camelCase (parts[0] + Title(rest...)).
func toCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}
	var b strings.Builder
	b.WriteString(parts[0])
	for _, p := range parts[1:] {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		b.WriteString(p[1:])
	}
	return b.String()
}

// camelKeys recursively rewrites map keys to camelCase.
//
// Slices are handled via reflection rather than a `case []any` type switch:
// a type switch on []any only matches that exact concrete type, not other
// slice types such as []map[string]any or []SomeStruct that handlers
// actually pass, so those would silently skip camel-casing entirely.
func camelKeys(obj any) any {
	if v, ok := obj.(map[string]any); ok {
		out := make(map[string]any, len(v))
		for k, val := range v {
			out[toCamel(k)] = camelKeys(val)
		}
		return out
	}

	rv := reflect.ValueOf(obj)
	if !rv.IsValid() {
		return obj
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		n := rv.Len()
		out := make([]any, n)
		for i := 0; i < n; i++ {
			out[i] = camelKeys(rv.Index(i).Interface())
		}
		return out
	case reflect.Map:
		out := make(map[string]any, rv.Len())
		for _, key := range rv.MapKeys() {
			out[toCamel(key.String())] = camelKeys(rv.MapIndex(key).Interface())
		}
		return out
	default:
		return obj
	}
}

// JSON writes the standard envelope to the Echo context.
// Pass data as map[string]any / []any to get automatic camelCasing,
// matching the legacy behaviour. Structs should carry camelCase json tags.
func JSON(c echo.Context, statusCode int, message string, data any, errVal any) error {
	var payload any
	if data != nil {
		payload = camelKeys(data)
	}
	return c.JSON(statusCode, Envelope{
		StatusCode: statusCode,
		Message:    message,
		Data:       payload,
		Error:      errVal,
	})
}

// OK is a convenience for 200 responses with data.
func OK(c echo.Context, message string, data any) error {
	return JSON(c, 200, message, data, nil)
}

// Err is a convenience for error responses.
func Err(c echo.Context, statusCode int, message string, errVal any) error {
	return JSON(c, statusCode, message, nil, errVal)
}
