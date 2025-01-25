# AST Instrumentation with go-instrument

## Instalasi
<!-- `go install github.com/nikolaydubina/go-instrument@1.7.0` -->
`go get -u github.com/ansrivas/fiberprometheus/v2`

## Query
Bisa pakai query ini
1. Durasi Lama Waktu Request per Route
`histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, route))`

2. Banyak HTTP OK dan BadRequest per Route
`sum(rate(http_requests_total{status="200"}[5m])) by (route)`

3. Gabungan status dan route
`sum(rate(http_requests_total[5m])) by (route, status)`

4. Sesuai dengan route
misal
```bash
http_requests_total{method="GET", route="/api/v1/users"} 
http_requests_total{method="POST", route="/api/v1/login"}
```


## Instrumentasi
`go-instrument -file main.go`

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