# FROM golang:1.23.4 AS builder
# LABEL maintainer="rbennum"
# WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download
# COPY . ./
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .
# RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# FROM alpine:latest
# WORKDIR /app
# COPY --from=builder /app/main .
# COPY --from=builder /app/.env .
# FROM debian:12-slim
# RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
# WORKDIR /root/
# COPY --from=builder /app/main .
# COPY --from=builder /app/.env .
# ENTRYPOINT [ "/app/main" ]
# CMD ["/bin/sh", "-c", "/root/main; tail -f /dev/null"]

# FROM golang:1.23.4-alpine AS build
# WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download
# COPY . .
# RUN CGO_ENABLED=0 GOOS=linux go build -o file-svc .

# FROM alpine:edge
# WORKDIR /app
# COPY --from=build /app/.env .
# COPY --from=build /app/file-svc .
# RUN apk --no-cache add ca-certificates tzdata
# ARG PORT=8080
# ARG GRPC_PORT=50051
# ENV PORT=${PORT}
# ENV GRPC_PORT=${GRPC_PORT}
# EXPOSE ${PORT} ${GRPC_PORT}
# ENTRYPOINT [ "/app/file-svc" ]

# multi-stage build
FROM golang:1.23.4 AS build
WORKDIR /app
COPY . .
RUN GOOS=linux go build -o file-svc .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /app/.env .
COPY --from=build /app/file-svc .
ARG MIGRATION_FILE_PATH
ENV MIGRATION_FILE_PATH=${MIGRATION_FILE_PATH}
COPY --from=build /app/${MIGRATION_FILE_PATH} ./${MIGRATION_FILE_PATH}
# RUN apt-get update && apt-get install -y strace
ENTRYPOINT [ "/app/file-svc" ]

# single-stage build
# FROM golang:1.23.4
# WORKDIR /app
# COPY . .
# RUN go build -o file_svc .
# ENTRYPOINT [ "/app/file_svc" ]