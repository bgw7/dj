package internal

type Location struct {
	ID      string
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode int    `json:"zipCode"`
}

type Venue struct {
	ID                  string
	PublicID            string
	Name                string
	SupportedEventTypes string
	Location            *Location
	Metadata
}

type VenueSearch struct {
	ClientID       *string `json:"clientId"`
	StartTimestamp *string `json:"startTimestamp"`
	City           *string `json:"city"`
	State          *string `json:"state"`
	ZipCode        *int    `json:"zipCode"`
	Limit          int     `json:"limit"`
	Offset         int     `json:"offset"`
}
