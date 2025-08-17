FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o myapp .

FROM alpine:latest

RUN apk add --no-cache sqlite-libs

COPY --from=builder /app/myapp /usr/local/bin/myapp

COPY ./eclipse.db /app/eclipse.db

WORKDIR /app

CMD ["myapp"]