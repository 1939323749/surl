FROM golang:1.20-alpine as builder

WORKDIR /app
ENV GOPROXY=https://goproxy.io,direct
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o shorturl

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/shorturl /shorturl

ENTRYPOINT ["/shorturl"]