FROM golang:1.24 as builder

ENV GOPATH=/
WORKDIR /app
COPY ./ /app

RUN go mod download && CGO_ENABLED=0 go build -o /keeper ./cmd/run/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /keeper /app/keeper
COPY ./config/dev.yaml /app/config/dev.yaml
COPY ./db/migrations /app/db/migrations

CMD ["/app/keeper"]