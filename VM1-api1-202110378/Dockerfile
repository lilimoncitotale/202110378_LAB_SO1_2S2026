FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o api1 main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api1 .
EXPOSE 5001
CMD ["./api1"]
