.PHONY: migrate register-admin build test run

migrate:
	go run cmd/migrator/main.go

register-admin:
	go run cmd/cli/main.go register-admin \
		--email admin@gmail.com --fullname Alex --phone 89111111111 --password qwerty1234 \
		--passport_series 1234 --passport_number 123456 --passport_issued_by unknown

build:
	go build -o main cmd/main/main.go
	go build -o migrate cmd/migrator/main.go
	go build -o cli cmd/cli/main.go

test:
	go test ./...

run:
	go run cmd/main/main.go
