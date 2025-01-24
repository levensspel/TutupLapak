package Entity

import "time"

type PaymentDetail struct {
	PurchaseID        string
	SellerID          string
	BankAccountName   string
	BankAccountHolder string
	SenderContactType int8
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
