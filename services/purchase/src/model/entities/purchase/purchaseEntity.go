package Entity

import "time"

type Purchase struct {
	PurchaseID          string
	SenderName          string
	SenderContactDetail string
	SenderContactType   int8
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
