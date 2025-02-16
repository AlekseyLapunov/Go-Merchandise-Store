FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go test -v -c -o e2e .

CMD ["./e2e"]