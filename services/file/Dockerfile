FROM golang:1.23.4 AS build
LABEL maintainer="rbennum"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o file-svc .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /app/.env .
COPY --from=build /app/file-svc .
COPY --from=build /app/src/database/migrations src/database/migrations
ENTRYPOINT [ "/app/file-svc" ]
