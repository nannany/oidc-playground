FROM golang:1.23 AS go
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/main .
CMD ["/app/main"]