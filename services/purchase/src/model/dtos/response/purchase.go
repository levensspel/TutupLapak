package response

import "time"

// PurchaseResponseDTO represents the response structure for a purchase.
type PurchaseResponseDTO struct {
	PurchaseId     string                `json:"purchaseId"`
	PurchasedItems []ProductItemDTO      `json:"purchasedItems"`
	TotalPrice     float64               `json:"totalPrice"`
	PaymentDetails []SellerBankDetailDTO `json:"paymentDetails"`
}

// PurchasedItemDTO represents the details of each purchased item.
type ProductItemDTO struct {
	ProductId        string    `json:"productId"`
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
	SellerId         string
}

// PaymentDetailDTO represents the details of payment for each seller.
type SellerBankDetailDTO struct {
	SellerId          string
	BankAccountName   string  `json:"bankAccountName"`
	BankAccountHolder string  `json:"bankAccountHolder"`
	BankAccountNumber string  `json:"bankAccountNumber"`
	TotalPrice        float64 `json:"totalPrice"`
}
