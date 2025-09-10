FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/go_social/main.go

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM golang:1.24-alpine

WORKDIR /app

COPY .env .env
COPY --from=builder /app/app .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

COPY ./internal/migration ./internal/migration

EXPOSE 3000


CMD source .env && migrate -database "postgres://go_social:go_social@db:5432/go_social?sslmode=disable"    -path ./internal/migration up && ./app

