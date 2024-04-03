# todo(averysmalldog): check that this new go version exists and builds correctly
FROM golang:1.21.1 AS builder
WORKDIR /go/src/github.com/averysmalldog/tesla-gen3wc-monitor/
# todo(averysmalldog): eliminate this fetch step
RUN go get -d -v github.com/averysmalldog/polly
COPY main.go .
RUN GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o app .

FROM scratch
COPY --from=builder /go/src/github.com/averysmalldog/tesla-gen3wc-monitor/app .
ENTRYPOINT [ "/app" ]