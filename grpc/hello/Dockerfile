FROM golang:1.14 AS builder

COPY . /workdir
WORKDIR /workdir
RUN go build -tags netgo -o /app

FROM alpine

COPY --from=builder /app /app

ENTRYPOINT ["/app"]
