package models

type SearchResult struct {
	Capsules []Capsule `json:"capsules"`
	Count    int       `json:"count"`
	Query    string    `json:"query"`
}
