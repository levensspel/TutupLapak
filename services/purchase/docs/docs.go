// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/purchase": {
            "post": {
                "description": "Pembeli akan memasukkan detail produk dan jumlah yang akan dibeli, kemudian mengembalikan daftar detail produk beserta dengan daftar detail bank dari masing-masing penjual",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Purchase"
                ],
                "summary": "Add items to the cart",
                "parameters": [
                    {
                        "description": "Cart Data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CartDto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "success response",
                        "schema": {
                            "$ref": "#/definitions/response.PurchaseResponseDTO"
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "request.CartDto": {
            "type": "object",
            "required": [
                "purchasedItems",
                "senderContactDetail",
                "senderContactType",
                "senderName"
            ],
            "properties": {
                "purchasedItems": {
                    "description": "array | minItems: 1",
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "$ref": "#/definitions/request.PurchasedItem"
                    }
                },
                "senderContactDetail": {
                    "description": "string | required | validates based on contact type",
                    "type": "string"
                },
                "senderContactType": {
                    "description": "string | required | enum: \"email\" / \"phone\"",
                    "type": "string",
                    "enum": [
                        "email",
                        "phone"
                    ]
                },
                "senderName": {
                    "description": "string | required | minLength: 4 | maxLength: 55",
                    "type": "string",
                    "maxLength": 55,
                    "minLength": 4
                }
            }
        },
        "request.PurchasedItem": {
            "type": "object",
            "required": [
                "productId",
                "qty"
            ],
            "properties": {
                "productId": {
                    "description": "string | required | should be a valid productId",
                    "type": "string"
                },
                "qty": {
                    "description": "number | required | min: 2",
                    "type": "integer",
                    "minimum": 2
                }
            }
        },
        "response.PaymentDetailDTO": {
            "type": "object",
            "properties": {
                "bankAccountHolder": {
                    "type": "string"
                },
                "bankAccountName": {
                    "type": "string"
                },
                "bankAccountNumber": {
                    "type": "string"
                },
                "totalPrice": {
                    "type": "number"
                }
            }
        },
        "response.PurchaseResponseDTO": {
            "type": "object",
            "properties": {
                "paymentDetails": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.PaymentDetailDTO"
                    }
                },
                "purchaseId": {
                    "type": "string"
                },
                "purchasedItems": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.PurchasedItemDTO"
                    }
                },
                "totalPrice": {
                    "type": "number"
                }
            }
        },
        "response.PurchasedItemDTO": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "fileId": {
                    "type": "string"
                },
                "fileThumbnailUri": {
                    "type": "string"
                },
                "fileUri": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "productId": {
                    "type": "string"
                },
                "qty": {
                    "type": "integer"
                },
                "sku": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
