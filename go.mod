module RWTAPI

go 1.23.2 // Align with toolchain and Dockerfile

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/gorilla/mux v1.8.1
	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c
)

require filippo.io/edwards25519 v1.1.0 // indirect
