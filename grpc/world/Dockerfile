FROM dictav/golang-tools:1.14 AS builder

COPY . /workdir
WORKDIR /workdir
RUN go build -o /app

FROM alpine

COPY --from=builder /app /app

ENTRYPOINT ["./app"]
