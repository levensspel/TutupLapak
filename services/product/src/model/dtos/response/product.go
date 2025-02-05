package response

import "time"

type Product struct {
	ProductId string `json:"productId"`
	Name      string `json:"Name"`
	Category  string `json:"category"`
	Qty       int    `json:"qty"`
	Price     int    `json:"price"`
	Sku       string `json:"sku"`
	FileId    string `json:"fileId"`
}

type ProductCreate struct {
	ProductId        string    `json:"productId"`
	Name             string    `json:"Name"`
	Category         string    `json:"category"`
	Qty              int       `json:"qty"`
	Price            int       `json:"price"`
	Sku              string    `json:"sku"`
	FileId           string    `json:"fileId"`
	FileUri          string    `json:"fileUri"`
	FileThumbnailUri string    `json:"fileThumbnailUri"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
