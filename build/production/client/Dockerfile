FROM golang:alpine AS builder

WORKDIR /build

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o main main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates curl tzdata

COPY --from=builder /build/main /app/main

RUN mkdir -p /app/.td && chmod 777 /app/.td

CMD ["/app/main"]