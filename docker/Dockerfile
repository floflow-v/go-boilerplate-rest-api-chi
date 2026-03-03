# Build
FROM golang:latest AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod .
RUN go mod download

COPY . .

RUN swag init --generalInfo cmd/go-rest-api-chi-example/main.go --parseDependency --parseInternal
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o ./go-rest-api-chi-example ./cmd/go-rest-api-chi-example

# Run
FROM gcr.io/distroless/static-debian12

WORKDIR /app
# Copy the built binary
COPY --from=builder /app/go-rest-api-chi-example /app/go-rest-api-chi-example
COPY --from=builder /app/docs /app/docs
ENTRYPOINT ["/app/go-rest-api-chi-example"]
