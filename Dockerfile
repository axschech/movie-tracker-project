FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN ls -la
RUN go build -o app -tags viper_bind_structs ./cmd/app/...

FROM alpine:latest

COPY --from=builder /app/app /app

CMD ["/app"]
