package models

type SearchBody struct {
	TypeSearch string `json:"type_search"`
	Term       string `json:"term"`
}
