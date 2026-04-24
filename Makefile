# =========================
# DATABASE
# =========================

DB_URL=postgres://postgres:postgres@localhost:5432/chainpulse?sslmode=disable

migUp:
	migrate -path=sql/migrations -database "$(DB_URL)" -verbose up

migDown:
	migrate -path=sql/migrations -database "$(DB_URL)" -verbose down

# =========================
# SQLC
# =========================

sqlc:
	sqlc generate

# =========================
# PROTO (BUF)
# =========================

proto:
	buf generate

proto-format:
	buf format -w

proto-lint:
	buf lint

proto-breaking:
	buf breaking --against '.git#branch=main'

# =========================
# CLEAN
# =========================

proto-clean:
	find . -name "*.pb.go" -delete
	find gen/openapiv2 -name "*.json" -delete

# =========================
# ALL
# =========================

generate: proto-format sqlc proto

.PHONY: migUp migDown sqlc proto proto-clean generate proto-format proto-lint proto-breaking