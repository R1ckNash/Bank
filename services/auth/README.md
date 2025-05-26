# Auth Service

```
auth/
├── internal/
│   ├── app/
│   │   ├── models/      // domain-models: User, Token etc
│   │   ├── repository/  // repo: auth_storage, etc
│   │   ├── server/      
│   │   ├── services/    // domain-models: User, Token etc
│   │   └── usecases/    // business-logic: user_auth/
│   │        └── user_auth.go
│   │     
│   └── middleware/      // errors, idempotency, etc
│      
├── api
├── cmd
├── pkg
└── go.mod

	/*
		-> [auth] -> [DB]
		     |
			  -> [Kafka] -> [Analytic-consumer] -> [Clickhouse] -> report
	*/

```