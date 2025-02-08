include .env

run:
	go run .

create_migration:
	@if [ -z "$(title)" ]; then \
		echo "error: please specify a title first"; \
		exit 1; \
	fi
	migrate create -ext=sql -dir=src/database/migrations -seq $(title)

migrate_up:
	migrate -path=database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

clean_uploads:
	find .uploads -type f -delete

doc_build:
	docker build -t file-service .

doc_run:
	docker run --rm -it -p 3003:3003 -p 5003:5003 --network demo-network file-service

.PHONY: run create_migration migrate_up migrate_down clean_uploads
