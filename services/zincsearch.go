package services

type ApiZincConector interface {
	search(body string) []Email
}
