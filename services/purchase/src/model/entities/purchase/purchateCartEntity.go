package Entity

import "time"

type PurchaseCart struct {
	PurchaseID string
	ProductID  string
	Quantity   int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
