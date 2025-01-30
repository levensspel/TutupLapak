package Entity

import "time"

type Purchase struct {
	PurchaseID          string
	SenderName          string
	SenderContactDetail string
	SenderContactType   string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
