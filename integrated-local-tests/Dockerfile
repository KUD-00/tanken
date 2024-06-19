FROM golang:1.22.3 as builder

WORKDIR /src

COPY integreted-local-tests/ /src/test/

WORKDIR /src/test

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /src/test/test ./main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /src/test/test /app/test

EXPOSE 50051

CMD ["/app/test"]