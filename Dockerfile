FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLE=0 go build -ldflags="-s -w" -o stresser .

FROM scratch

COPY --from=builder /app/stresser /usr/bin/stresser
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["stresser"]
