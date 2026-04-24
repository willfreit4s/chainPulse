migUp:
	migrate -path=sql/migrations -database "postgres://postgres:postgres@localhost:5432/chainpulse?sslmode=disable" -verbose up

migDown:
	migrate -path=sql/migrations -database "postgres://postgres:postgres@localhost:5432/chainpulse?sslmode=disable" -verbose down

.PHONY: migrateUp migrateDown createMigration

generate:
	sqlc generate
