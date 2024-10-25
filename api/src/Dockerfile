FROM --platform=linux/amd64 golang:1.23-alpine3.20 AS builder

ENV CGO_ENABLED=1
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite \
    sqlite-dev \
    sqlite-libs

COPY . /go/src/com.github/w-k-s/glassdoor-hr-review-detector
WORKDIR /go/src/com.github/w-k-s/glassdoor-hr-review-detector

RUN go get ./...
# Generating static binary: https://7thzero.com/blog/golang-w-sqlite3-docker-scratch-image
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags '-linkmode external -extldflags "-static"' -o app com.github/w-k-s/glassdoor-hr-review-detector/cmd/server

FROM alpine:latest

RUN apk add --no-cache \
    ca-certificates

WORKDIR /root/
RUN mkdir -p /root/.glassdoor-hr-review-detector/migrations

COPY --from=builder /go/src/com.github/w-k-s/glassdoor-hr-review-detector/app .
COPY --from=builder /go/src/com.github/w-k-s/glassdoor-hr-review-detector/migrations /root/.glassdoor-hr-review-detector/migrations

ENTRYPOINT ["./app"]