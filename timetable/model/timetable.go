package model

type Period struct {
	Type    string `json:"type"`
	Subject string `json:"subject"`
	Book    string `json:"book"`
	Name    string `json:"name"`
}

type Day struct {
	Name        string   `json:"name"`
	UniformType string   `json:"uniform"`
	Periods     []Period `json:"periods"`
}

type Timetable struct {
	Schoolcode string `json:"school"`
	Grade      string `json:"grade"`
	Division   string `json:"division"`
	Days       []Day  `json:"days"`
}
