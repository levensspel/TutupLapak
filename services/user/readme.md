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

### NOTE: 
it is important to put the build inside of the .build folder
to ensure the gitignore caught up with the files.
