package problem

// Request bodies ported from schemas/problem.py (CamelModel -> camelCase JSON).

type createProblemRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Hints       []string `json:"hints"`
	Difficulty  string   `json:"difficulty"`
	SectionID   *string  `json:"sectionId"`
}

type updateProblemRequest struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Hints       *[]string `json:"hints"`
	Difficulty  *string   `json:"difficulty"`
}

type addTopicRequest struct {
	TagName []string `json:"tagName"`
}

type updateTopicRequest struct {
	TagName []string `json:"tagName"`
}

type addTemplateRequest struct {
	FunctionBody string `json:"functionBody"`
	DefaultCode  string `json:"defaultCode"`
	Language     string `json:"language"`
}

type updateTemplateRequest struct {
	FunctionBody *string `json:"functionBody"`
	DefaultCode  *string `json:"defaultCode"`
	Language     *string `json:"language"`
}

type addTestCaseRequest struct {
	TestType string `json:"testType"`
	Input    string `json:"input"`
	Output   string `json:"output"`
}

type updateTestCaseRequest struct {
	TestType *string `json:"testType"`
	Input    *string `json:"input"`
	Output   *string `json:"output"`
}

type addSolutionRequest struct {
	Solution      string  `json:"solution"`
	VideoSolution *string `json:"videoSolution"`
}

type updateSolutionRequest struct {
	Solution      *string `json:"solution"`
	VideoSolution *string `json:"videoSolution"`
}

type searchProblemRequest struct {
	Query string `json:"query"`
}
