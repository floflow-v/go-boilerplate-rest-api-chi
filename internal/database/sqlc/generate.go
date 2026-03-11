package sqlc

//go:generate mockgen -destination=../../mocks/mock_querier.go -package=mocks go-rest-api-chi-example/internal/database/sqlc Querier
