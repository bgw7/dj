package internal

type Track struct {
	ID          int    `json:"id"`
	Url         string `json:"url"`
	Filename    string `json:"filename"`
	VoteCount   int    `json:"voteCount"`
	CreatedAt   string `json:"createdAt"`
	HasPlayed   bool   `json:"hasPlayed"`
	CreatedBy   string `json:"createdBy"`
	CreatedWith string `json:"createdWith"`
}

type Vote struct {
	Filename string `json:"filename"`
	Url      string `json:"url"`
	VoterID  string `json:"voterId"`
}
