FROM golang:1.25-alpine AS modules
WORKDIR /modules

COPY go.mod ./
RUN go mod download

FROM golang:1.25-alpine AS builder
WORKDIR /app

# Certs for HTTPS 
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Modules cache
COPY --from=modules /go/pkg/ /go/pkg
COPY --from=modules /go/src /go/src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app cmd/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/app /app

EXPOSE 8080

CMD ["/app"]
