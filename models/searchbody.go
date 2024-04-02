package models

type SearchBody struct {
	TypeSearch string `json:"type_search"`
	Term       string `json:"term"`
	From       int    `json:"from"`
	MaxResults int    `json:"max_results"`
}
