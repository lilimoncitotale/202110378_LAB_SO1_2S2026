WORKDIR /app
COPY main.go .
RUN go build -o api2 main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api2 .
EXPOSE 5002
CMD ["./api2"]
