run:
	go run cmd/api/main.go

db-generate:
	xo schema "mysql://user:pass@127.0.0.1:3306/accounting_system" -o internal/models --include users --include accounts --include journal_vouchers --include journal_entries
