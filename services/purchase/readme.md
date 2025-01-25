# Cache
Dengan graph-io/ristretto:2.3.0

## Installation
`go get github.com/dgraph-io/ristretto/v2`

# Api Documentation

## Installation
```bash
go get -u github.com/swaggo/swag/cmd/swag
go install github.com/swaggo/swag/cmd/swag@1.1.1
```

## Add Documentation
```bash
swag init
```

## Add Swagger
```bash
go get -u github.com/gofiber/swagger
```

## Menambahkan endpoint baru
```go
// Contoh dibawah ini

// Cart godoc 
// @Summary Add items to the cart
// @Description Endpoint to add purchased items into the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body request.CartDto true "Cart Data"
// @Success 200 {object} map[string]interface{} "success response"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Router /v1/cart [post]

```

Kemudian jalankan `swag init`

# Menjalankan Aplikasi
1. Copy-paste terlebih dulu `.env.example` dan ganti menjadi `.env`
2. Masukan kredensial pada local development
3. `go run main.go`