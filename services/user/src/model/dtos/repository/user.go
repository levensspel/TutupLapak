package repository

type AuthByEmail struct {
	UserId       string
	HashPassword string
	Phone        string
}

type AuthByPhone struct {
	UserId       string
	HashPassword string
	Email        string
}

type User struct {
	Email             *string
	Phone             *string
	FileId            *string
	FileUri           *string
	FileThumbnailUri  *string
	BankAccountName   *string
	BankAccountHolder *string
	BankAccountNumber *string
}

type UpdateUser struct {
	FileId            *string
	FileUri           *string
	FileThumbnailUri  *string
	BankAccountName   *string
	BankAccountHolder *string
	BankAccountNumber *string
}
