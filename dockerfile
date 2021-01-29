FROM golang:1.15.6 AS builder
WORKDIR /go/src/github.com/averysmalldog/tesla-gen3wc-monitor/
RUN go get -d -v github.com/averysmalldog/polly  
COPY main.go .
RUN GOOS=linux go build -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/averysmalldog/tesla-gen3wc-monitor/app .
CMD ["./app"]