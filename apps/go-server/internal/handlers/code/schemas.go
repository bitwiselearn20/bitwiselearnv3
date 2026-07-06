package code

type runCodeRequest struct {
	Code      string `json:"code"`
	Language  string `json:"language"`
	ProblemID string `json:"problemId"`
}

type compileCodeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Stdin    string `json:"stdin"`
}

type submitCodeRequest struct {
	Code      string `json:"code"`
	Language  string `json:"language"`
	ProblemID string `json:"problemId"`
}
