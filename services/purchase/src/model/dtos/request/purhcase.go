package request

/* Purchase to put products in an cart
{
  "purchasedItems": [ // array | minItems: 1
    {
      "productId": "", // string | should a valid productId
      "qty": 1, // number | min: 2
    },
  ],
  "senderName": "", // string | required | minLength: 4 | maxLength: 55
  "senderContactType": "", // string | required | enum of "email" / "phone"
  "senderContactDetail": "", // string | required | if "phone" then validates the phone number | if "email" then validates email
}
*/

// PurchasedItem defines the structure for an item in the cart
type PurchasedItem struct {
	ProductId string `json:"productId" validate:"required,uuid4"` // Ensure productId is a valid UUID
	Qty       int    `json:"qty" validate:"required,min=2"`       // Minimum quantity is 2
}

// CartDto represents the purchase request
type CartDto struct {
	PurchasedItems      []PurchasedItem `json:"purchasedItems" validate:"required,min=1,dive"`                 // Array must have at least 1 item
	SenderName          string          `json:"senderName" validate:"required,min=4,max=55"`                   // Sender name constraints
	SenderContactType   string          `json:"senderContactType" validate:"required,oneof=email phone"`       // Must be "email" or "phone"
	SenderContactDetail string          `json:"senderContactDetail" validate:"required,sender_email_or_phone"` // Conditional validation
}
