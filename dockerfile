FROM golang:latest AS builder
WORKDIR /go/src/github.com/go-backend-test
COPY . .
RUN go mod download
RUN go build -o main .

FROM golang:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/go-backend-test/main .
COPY --from=builder /go/src/github.com/go-backend-test/.env .
EXPOSE 8080
CMD ["./main"]