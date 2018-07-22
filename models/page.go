package models

type Page struct {
	Page         uint `json:"page"`
	TotalResults uint `json:"totalResults"`
	TotalPages   uint `json:"totalPages"`
}
