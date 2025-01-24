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

type PurchasedItem struct {
	ProductId string `json:"productId" binding:"required"` // string | required | should be a valid productId
	Qty       int    `json:"qty" binding:"required,min=2"` // number | required | min: 2
}

type CartDto struct {
	PurchasedItems      []PurchasedItem `json:"purchasedItems" binding:"required,min=1"`                // array | minItems: 1
	SenderName          string          `json:"senderName" binding:"required,min=4,max=55"`             // string | required | minLength: 4 | maxLength: 55
	SenderContactType   string          `json:"senderContactType" binding:"required,oneof=email phone"` // string | required | enum: "email" / "phone"
	SenderContactDetail string          `json:"senderContactDetail" binding:"required"`                 // string | required | validates based on contact type
}
