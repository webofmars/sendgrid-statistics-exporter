FROM golang:1.14.9-alpine AS builder

RUN apk add --update --no-cache git ca-certificates

WORKDIR /opt

COPY ./ ./

ENV CGO_ENABLED=0

RUN go mod download
RUN go build -o /opt/sendgrid-exporter /opt/cmd/sendgrid-exporter/main.go

FROM scratch

COPY --from=builder /opt/sendgrid-exporter /opt/sendgrid-exporter
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/opt/sendgrid-exporter"]