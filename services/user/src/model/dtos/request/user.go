package request

type AuthByEmailRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type AuthByPhoneRequest struct {
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type LinkEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}
