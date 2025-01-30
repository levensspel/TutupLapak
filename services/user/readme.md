# TutupLapak User Service

### Build & Run the App
Ensure make an `.env` file with values filling the keys from the `.env.example` file

To build:
```bash
go build -o .build/<name-of-build.extension?>
```

To run the binary:
```bash
./.build/<name-of-build.extension?>
```

To quick run without build:
```bash
go run main.go
```

# gRPC
## Installation
```bash
# 1. Install protobuf compiler, 
# Ubuntu / iOS
sudo apt install protobuf-compiler
# Windows, buka powershell Administrator dulu
choco install protoc

# 2. Pasang auto-generated go proto
go install google.golang.org/protobuf/cmd/protoc-gen-go@1.36.4
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@1.5.1

# 3. Pasang depedensi untuk service golang sendiri
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u google.golang.org/protobuf
```
## Setup
1. Buat dulu kontrak proto nya, bisa copas .proto di ./services/user/proto
2. Jalankan protoc untuk generate grpc golang services
```bash
protoc --go_out=plugins=grpc:. ./proto/*.proto
#atau
protoc --go_out=. --go-grpc_out=. ./proto/*.proto
```
3. Setelah di generate. Pakai _Service yang ada di ..._grpc.pb.go
4. Pada `server.go`, siapkan grpc server dan daftarkan terlebih dulu --memang belum ada implementasi controller untuk service grpc, lebih baik disiapkan--
```go
go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		puc := do.MustInvoke[*protoUserController.ProtoUserController](di.Injector)
		user.RegisterUserServiceServer(grpcServer, _)

		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
```
5. Untuk membuat controller, implementasinya serupa tanpa grpc. Hanya saja ada perbedaan pada New___Inject()
```go
func NewUserControllerInject(i do.Injector) (*ProtoUserController, error) {
    ...
    return &ProtoUserController{...}, nil
}
```
6. Pada `./di/injector.go`
```go
do.Provide[*protoUserController.ProtoUserController](Injector, protoUserController.NewUserControllerInject)
```
7. Jalankan aplikasi seperti biasa, bisa dengan 
```bash
go run main.go
```
8. Lakukan bugfix apabila ada masalah atau aplikasi tidak mau jalan

## Test
1. Dengan postman, bisa lakukan `New > gRPC > ...`
2. Masukan URL, misalnya
```go
localhost:50051
```
3. Pilih metode, dan pilih upload .proto
4. Upload proto yang sudah kamu buat sebelumnya, misal `./proto/user_service.proto` yang kemudian akan kita jadikan sebuah kontrak proto
5. Klik pada metode lagi, lalu pilih salah satu method yang ada
6. Lakukan pengujian dengan menekan tombol "Invoke"
7. Kalau mau ada payload, pilih pada tab "message"
8. Buat json seperti biasa, misal
```json
{
  "userIds": ["user1", "user2", "user3"]
}
```

## Prometheus for Grpc
1. Install
```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
go get github.com/grpc-ecosystem/go-grpc-prometheus
```
2. 

### NOTE: 
it is important to put the build inside of the .build folder
to ensure the gitignore caught up with the files.

# Redis
## Setup
#### 1. Ensure you have Redis installed or you can use the Redis with container.


Example to run the Redis using existing Docker Compose  in this repository:
```
cd ./middleware/cacheserver
docker compose up -d 
```

#### 2. Check if Redis container run properly
Enter the Redis container:
```
docker compose exec -it redis bash
```
Run Redis CLI inside the redis:
```
redis-cli
```
Ping the Redis server to see whether it answers "PONG"!
```
ping
```

#### 3. Ensure the Redis environtment variables are all set properly in `.env` file.

#### 4. Run the app like usual
```bash
go run main.go
```

## Test The Redis
After running the app and testing related API endpoints, check whether the Redis really stores the cache data.
Enter Redis CLI:
```
docker compose exec -it redis bash
redis-cli
```
Do any Redis syntax you want. 
Example to get value from key:
```
GET key_name
```

## Add Interfaces
1. Open `./services/user/src/cache/redis`
2. Add the new interface to `CacheClientInterface`
3. Add the function implementation to it and ensure the functon exposed to the consumer