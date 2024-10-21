FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY facts.txt ./
COPY main.go ./

RUN go build -o fastfact .

FROM alpine:latest

COPY --from=builder /app/fastfact /usr/local/bin/fastfact

EXPOSE 8080

CMD ["fastfact"]

