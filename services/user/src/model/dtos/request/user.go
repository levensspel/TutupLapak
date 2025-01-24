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

type LinkPhoneRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type UpdateUserProfileRequest struct {
	FileId            string `json:"fileId"`
	BankAccountName   string `json:"bankAccountName" validate:"required,min=4,max=32"`
	BankAccountHolder string `json:"bankAccountHolder" validate:"required,min=4,max=32"`
	BankAccountNumber string `json:"bankAccountNumber" validate:"required,min=4,max=32"`
}
