package internal

type PagedResponse[T any] struct {
	Offset  int `json:"offset"`
	Results []T `json:"results"`
	Total   int `json:"total"`
}
