package internal

import "time"

type Track struct {
	ID          int       `json:"id"`
	Url         string    `json:"url"`
	Filename    string    `json:"filename"`
	VoteCount   int       `json:"voteCount"`
	CreatedAt   time.Time `json:"createdAt"`
	HasPlayed   bool      `json:"hasPlayed"`
	CreatedBy   string    `json:"createdBy"`
	CreatedWith string    `json:"createdWith"`
}

type VoteRequest struct {
	ID      int    // Unique Vote ID for delete operations
	TrackID int    // Track that the vote is associated with
	VoterID string // Identifier for the voter
}

type Vote struct {
	Filename string `json:"filename"`
	Url      string `json:"url"`
	VoterID  string `json:"voterId"`
}

type DownloadRequest struct {
	URL string `json:"url"`
}
