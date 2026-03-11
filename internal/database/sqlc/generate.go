package sqlc

//go:generate mockgen -destination=../../mocks/mock_querier.go -package=mocks go-boilerplate-rest-api-chi/internal/database/sqlc Querier
