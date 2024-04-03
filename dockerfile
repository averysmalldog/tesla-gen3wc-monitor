FROM golang:1.21.1 AS builder

COPY . /go/src/github.com/averysmalldog/tesla-gen3wc-monitor/
WORKDIR /go/src/github.com/averysmalldog/tesla-gen3wc-monitor/
RUN go mod download

RUN GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o app .

FROM scratch
COPY --from=builder /go/src/github.com/averysmalldog/tesla-gen3wc-monitor/app .
ENTRYPOINT [ "/app" ]