FROM golang:1.19.9-alpine3.18 AS builder

RUN apk add --no-cache build-base

WORKDIR /book_service

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN go build -o booktastic_book_service .

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /book_service

COPY --from=builder /book_service/booktastic_book_service .

RUN chmod +x booktastic_book_service

EXPOSE 8080

CMD ["./booktastic_book_service"]
