package request

type ProductCreate struct {
	Name     *string `validate:"required,min=4,max=32"`
	UserId   string
	Category *string `validate:"required"`
	Qty      *int    `validate:"required,min=1"`
	Price    *int    `validate:"required,min=100"`
	Sku      *string `validate:"required,min=0"`
	FileId   *string `validate:"required"`
}
