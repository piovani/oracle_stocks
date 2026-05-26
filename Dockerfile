FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /oracle_stocks ./cmd/api

FROM alpine:3.22

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /oracle_stocks .

EXPOSE 8080

ENTRYPOINT ["./oracle_stocks"]
