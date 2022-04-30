FROM golang:latest AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src
COPY ./src /src
RUN go build -o /build/app  .

FROM alpine:latest
COPY ./env /env
COPY --from=builder /build/app .
ENTRYPOINT ["/app"]

