module auth

go 1.24

require (
	github.com/R1ckNash/Bank/pkg v0.0.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-chi/render v1.0.3
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/jackc/pgerrcode v0.0.0-20240316143900-6e2875d9b438
	github.com/jackc/pgx/v5 v5.7.5
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.37.0
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/georgysavva/scany/v2 v2.1.4 // indirect
	github.com/golang-migrate/migrate/v4 v4.18.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/vgarvardt/pgx-google-uuid/v5 v5.6.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

replace github.com/R1ckNash/Bank/pkg => ../../pkg
