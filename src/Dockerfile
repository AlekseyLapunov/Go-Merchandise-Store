FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o merch_app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/merch_app .
EXPOSE 8080
CMD ["./merch_app"]
