FROM golang:alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev git

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum* ./
RUN go mod download

CMD ["air", "-c", "build/development/client/.air.toml"]