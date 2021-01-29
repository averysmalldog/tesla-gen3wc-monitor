FROM golang:1.15.6 AS builder
WORKDIR /go/src/github.com/averysmalldog/tesla-gen3wc-monitor/
RUN go get -d -v github.com/averysmalldog/polly 
COPY main.go .
RUN GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o app .

FROM scratch
COPY --from=builder /go/src/github.com/averysmalldog/tesla-gen3wc-monitor/app .
ENTRYPOINT [ "/app" ]