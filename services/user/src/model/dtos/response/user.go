package response

import "database/sql"

type AuthResponse struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	Token string `json:"token"`
}

type UserResponse struct {
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	FileId            string `json:"fileId"`
	FileUri           string `json:"fileUri"`
	FileThumbnailUri  string `json:"fileThumbnailUri"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountHolder string `json:"bankAccountHolder"`
	BankAccountNumber string `json:"bankAccountNumber"`
}

type UserWithIdResponse struct {
	UserId            string `json:"userId"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	FileId            string `json:"fileId"`
	FileUri           string `json:"fileUri"`
	FileThumbnailUri  string `json:"fileThumbnailUri"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountHolder string `json:"bankAccountHolder"`
	BankAccountNumber string `json:"bankAccountNumber"`
}

type UserResponseSqlNulString struct {
	Email             sql.NullString `json:"email"`
	Phone             sql.NullString `json:"phone"`
	FileId            sql.NullString `json:"fileId"`
	FileUri           sql.NullString `json:"fileUri"`
	FileThumbnailUri  sql.NullString `json:"fileThumbnailUri"`
	BankAccountName   sql.NullString `json:"bankAccountName"`
	BankAccountHolder sql.NullString `json:"bankAccountHolder"`
	BankAccountNumber sql.NullString `json:"bankAccountNumber"`
}

type UserWithIdResponseSqlNulString struct {
	UserId            sql.NullString `json:"userId"`
	Email             sql.NullString `json:"email"`
	Phone             sql.NullString `json:"phone"`
	FileId            sql.NullString `json:"fileId"`
	FileUri           sql.NullString `json:"fileUri"`
	FileThumbnailUri  sql.NullString `json:"fileThumbnailUri"`
	BankAccountName   sql.NullString `json:"bankAccountName"`
	BankAccountHolder sql.NullString `json:"bankAccountHolder"`
	BankAccountNumber sql.NullString `json:"bankAccountNumber"`
}
