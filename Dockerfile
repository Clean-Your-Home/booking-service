FROM golang:1.22-alpine AS modules

WORKDIR /modules

COPY go.mod ./
RUN go mod download

FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app cmd/main.go

FROM scratch

COPY --from=builder /bin/app /app

EXPOSE 8080

CMD ["/app"]
