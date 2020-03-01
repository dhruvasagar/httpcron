FROM golang:1.12.7 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o httpcron .

FROM alpine:3.9
COPY .env .
COPY --from=builder /app/httpcron .
EXPOSE 9000
CMD ["./httpcron"]
