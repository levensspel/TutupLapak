package entity

import "time"

type Product struct {
	Id        string
	UserId    string
	Name      string
	Category  string
	Qty       int
	Price     int
	Sku       string
	FileId    string
	CreatedAt time.Time
	UpdateAt  time.Time
}
