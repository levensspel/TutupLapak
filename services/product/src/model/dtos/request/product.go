package request

type ProductCreate struct {
	Name     *string `validate:"required,min=4,max=32"`
	UserId   string
	Category *string `validate:"required,category_product"`
	Qty      *int    `validate:"required,min=1"`
	Price    *int    `validate:"required,min=100"`
	Sku      *string `validate:"required,min=0"`
	FileId   *string `validate:"required"`
}

type ProductUpdate struct {
	Id       string
	Name     *string `validate:"required,min=4,max=32"`
	UserId   string
	Category *string `validate:"required,category_product"`
	Qty      *int    `validate:"required,min=1"`
	Price    *int    `validate:"required,min=100"`
	Sku      *string `validate:"required,min=0"`
	FileId   *string `validate:"required"`
}

type ProductFilter struct {
	Limit     int
	Offset    int
	ProductId string `validate:"category_product"`
	Sku       string
	Category  string
	SortBy    string `validate:"category_search"`
}
