FROM golang:1.23-alpine AS dev

RUN go install github.com/air-verse/air@v1.61.7

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/services/progress

CMD ["air", "-c", ".air.toml"]
