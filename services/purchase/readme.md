
# gRPC Client
## Instalasi
1. Jalankan perintah dibawah dalam terminal
```bash
# 1. Install protobuf compiler, 
# Ubuntu / iOS
sudo apt install protobuf-compiler
# Windows, buka powershell Administrator dulu
choco install protoc

# 2. Pasang auto-generated go proto
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 3. Pasang depedensi untuk service golang sendiri
go get -u github.com/golang/protobuf/protoc-gen-go 
#atau 
go get -u google.golang.org/protobuf

```
## Setup
1. Ambil dulu `.proto` dan letakan dalam satu direktori `./proto/`
2. Jalankan command berikut
```bash
protoc --go_out=. --go-grpc_out=. ./proto/*.proto
```
3. Setelah di generate. Pakai _Service yang ada di ..._grpc.pb.go
4. Karena kita akan bangun client, buka `./src/http/server.go`
5. 

# Mengambil log dari docker-container
```bash
# Cek ContainerID atau ContainerName
docker ps
# Ambil ID dan ganti container_id dengan ID container yang ada
docker cp [container_id]:/root/logs/app/log ./local-log
```
# Prometheus dan Grafana

Untuk melihat dashboard pada Grafana, kita harus setup terlebih dulu Prometheus yang memungkinkan mengambil-scrap- data dari service

## Instalasi
<!-- `go install github.com/nikolaydubina/go-instrument@1.7.0` -->
1. Go Get fiberprometheus
```bash
go get -u github.com/ansrivas/fiberprometheus/v2
```
2. Pasang terlebih dulu middleware fiberprometeus di `server.go`
```go
prometheus := fiberprometheus.New("ur-service-name")
prometheus.RegisterAt(app, "/metrics")
app.Use(prometheus.Middleware)
```
3. Buat `prometheus.yml` pada direktori `./src/config`
```yml
# Copy paste konfigurasi prometheus
scrape_configs:
  - job_name: "purhcase-service"
    static_configs:
      - targets: ["tutuplapak-purchase-service:8080"]
```
4. Pada `docker-compose.yml`
```yml
  prometheus:
    image: prom/prometheus
    volumes:
      - ./src/config/prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
```
5. Jalankan docker-compose dan buka pada halaman
```
http://localhost:9095/
```
6. Buka tab status>target health dan pastikan service kamu OK

7. Kunjungi laman ini https://grafana.com/grafana/dashboards/14331-fiber-framework-processes/ dan download template dashboard https://grafana.com/api/dashboards/14331/revisions/3/download

8. Akses `http://localhost:3000/` dan buka pada laman dashboard

9. Import dashboard, dan upload .json yang sudah kita download sebelumnya
! Jika ditanya data-akses prometheus buat baru

10. arahkan url prometheus pada `http://[ur_prometheus_service_container]:8080`

11. lakukan test and save

12. kembali pada laman dashboard, lakukan lahkah 9 

## Query
Bisa pakai query ini
1. Durasi Lama Waktu Request per Route
```bash
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, route))
```

2. Banyak HTTP OK dan BadRequest per Route
```bash
sum(rate(http_requests_total{status="200"}[5m])) by (route)
```

3. Gabungan status dan route
```bash
sum(rate(http_requests_total[5m])) by (route, status)
```

4. Sesuai dengan route
misal
```bash
http_requests_total{method="GET", route="/api/v1/users"} 
http_requests_total{method="POST", route="/api/v1/login"}
```

# Cache
Dengan graph-io/ristretto:2.3.0

## Installation
```bash
go get github.com/dgraph-io/ristretto/v2
```

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
func foo(c *fiber.Ctx) error {...}
```

Kemudian jalankan 
```bash
swag init
```

# Menjalankan Aplikasi
1. Copy-paste terlebih dulu `.env.example` dan ganti menjadi `.env`
2. Masukan kredensial pada local development
3. `go run main.go`