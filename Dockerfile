FROM golang:1.21.1 as builder

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest

COPY --from=builder /app/app /app

EXPOSE 8080

CMD ["/app"]
