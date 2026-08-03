FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o wow-exporter ./cmd/exporter


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/wow-exporter .

COPY config.yaml .

EXPOSE 9105

CMD ["./wow-exporter"]