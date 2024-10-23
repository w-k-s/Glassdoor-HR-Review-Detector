FROM golang:1.23 as builder

COPY . /go/src/com.github/w-k-s/glassdoor-hr-review-detector

WORKDIR /go/src/com.github/w-k-s/glassdoor-hr-review-detector

RUN go get ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app com.github/w-k-s/glassdoor-hr-review-detector/cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

RUN mkdir -p /root/.glassdoor-hr-review-detector/migrations.d
COPY --from=builder /go/src/com.github/w-k-s/glassdoor-hr-review-detector/app .
COPY --from=builder /go/src/com.github/w-k-s/glassdoor-hr-review-detector/migrations .glassdoor-hr-review-detector/migrations.d

ENTRYPOINT ["./app"]