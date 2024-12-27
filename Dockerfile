# Build stage
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY .env.dev .env
COPY start.sh .
#COPY wait-for.sh .
COPY internal/db/migrations ./internal/db/migrations

EXPOSE 8080 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
