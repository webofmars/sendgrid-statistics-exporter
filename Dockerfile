FROM golang:latest

RUN go get gitlab.com/countsheep123/sendgrid-exporter/...

ENTRYPOINT $GOPATH/bin/sendgrid-exporter
