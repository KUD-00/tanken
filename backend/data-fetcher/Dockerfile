FROM golang:1.22.2 as builder

WORKDIR /src

COPY data-fetcher/ /src/data-fetcher/
COPY common/ /src/common/

WORKDIR /src/data-fetcher

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /src/data-fetcher/data-fetcher ./main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /src/data-fetcher/data-fetcher /app/data-fetcher

EXPOSE 50051

CMD ["/app/data-fetcher"]
