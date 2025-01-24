package response

import "time"

// PurchaseResponseDTO represents the response structure for a purchase.
type PurchaseResponseDTO struct {
	PurchaseID     string             `json:"purchaseId"`
	PurchasedItems []PurchasedItemDTO `json:"purchasedItems"`
	TotalPrice     float64            `json:"totalPrice"`
	PaymentDetails []PaymentDetailDTO `json:"paymentDetails"`
}

// PurchasedItemDTO represents the details of each purchased item.
type PurchasedItemDTO struct {
	ProductID        string    `json:"productId"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	Qty              int       `json:"qty"`
	Price            float64   `json:"price"`
	SKU              string    `json:"sku"`
	FileID           string    `json:"fileId"`
	FileURI          string    `json:"fileUri"`
	FileThumbnailURI string    `json:"fileThumbnailUri"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// PaymentDetailDTO represents the details of payment for each seller.
type PaymentDetailDTO struct {
	BankAccountName   string  `json:"bankAccountName"`
	BankAccountHolder string  `json:"bankAccountHolder"`
	BankAccountNumber string  `json:"bankAccountNumber"`
	TotalPrice        float64 `json:"totalPrice"`
}
