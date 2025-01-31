FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dbserver .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/dbserver .

ENTRYPOINT [ "/app/dbserver" ]