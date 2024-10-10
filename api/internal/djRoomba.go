package internal

type Track struct {
	ID        int     `json:"id"`
	Url       string  `json:"url"`
	Filename  *string `json:"filename"`
	VoteCount int     `json:"voteCount"`
	HasPlayed bool    `json:"hasPlayed"`
	CreatedBy string  `json:"createdBy"`
}

type Vote struct {
	TrackID int    `json:"trackId"`
	UserID  string `json:"userId"`
}
