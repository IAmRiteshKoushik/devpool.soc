package models

type IssueAction struct {
	ParticipantUsername string `json:"github_username"`
	Url                 string `json:"url"`
	Claim               bool   `json:"claimed"`
}

type BountyAction struct {
	ParticipantUsername string `json:"github_username"`
	Amount              int    `json:"amount"`
	Url                 string `json:"url"`
	Action              string `json:"action"`
}

type Solution struct {
	ParticipantUsername string `json:"github_username"`
	Url                 string `json:"pull_request_url"`
	Merged              bool   `json:"merged"`
}
