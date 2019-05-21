FROM golang:1.12-alpine3.9 AS builder
WORKDIR src/github.com/Dimitriy14/shortener
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/shortener .

FROM alpine:3.7
WORKDIR /root/
COPY --from=builder /go/src/github.com/Dimitriy14/shortener/config/config.json ./config.json
COPY --from=builder /go/src/github.com/Dimitriy14/shortener/static ./static
COPY --from=builder /go/bin/shortener .
CMD ["./shortener"]
