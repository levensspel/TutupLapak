package response

type AuthResponse struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	Token string `json:"token"`
}
