# Stage 1: Build
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o vrcomandaapi .

# Stage 2: Tempo de execucao
FROM alpine:3.20

RUN apk add --no-cache ca-certificates sqlite-libs tzdata

WORKDIR /app

RUN mkdir -p /app/data

COPY --from=builder /app/vrcomandaapi .

EXPOSE 28232

CMD ["./vrcomandaapi"]
