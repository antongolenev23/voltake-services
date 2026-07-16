package dto

type MeResponse struct {
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}
