# Build stage (1st stage)
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
# main is the name of the output binary file

# Run stage (final run stage)
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]