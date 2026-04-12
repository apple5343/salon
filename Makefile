migrate:
	migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path ./migrations up

register-admin:
	go run cmd/cli/main.go register-admin \
		--email admin@gmail.com --fullname Alex --phone 89111111111 --password qwerty1234 \
		--passport_series 1234 --passport_number 123456 --passport_issued_by unknown