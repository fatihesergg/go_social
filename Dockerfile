FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/go_social/main.go

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM golang:1.24-alpine

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

COPY ./internal/migration ./internal/migration
COPY .env ./.env

EXPOSE 3000
CMD ["sh", "-c", "migrate -path ./internal/migration -database \"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable\" up || exit 1 && ./app"]
