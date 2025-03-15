package internal

type Track struct {
	ID          int     `json:"id"`
	Url         string  `json:"url"`
	Filename    *string `json:"filename"`
	VoteCount   int     `json:"voteCount"`
	HasPlayed   bool    `json:"hasPlayed"`
	CreatedBy   string  `json:"createdBy"`
	CreatedAt   string  `json:"createdAt"`
	CreatedWith string  `json:"createdWith"`
}

type Vote struct {
	Filename string `json:"filename"`
	Url      string `json:"url"`
	VoterID  string `json:"voterId"`
}
