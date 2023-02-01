FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o fffff-frontend .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/fffff-frontend .

ENTRYPOINT ["./fffff-frontend"]
