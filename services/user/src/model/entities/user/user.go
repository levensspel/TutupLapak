package Entity

import "time"

type User struct {
	Id                string
	Email             string
	PasswordHash      string
	Phone             string
	FileId            string
	FileUri           string
	FileThumbnailUri  string
	BankAccountName   string
	BankAccountHolder string
	BankAccountNumber string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
