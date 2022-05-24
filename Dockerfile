FROM golang:latest AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://goproxy.io,direct \
    GOOS=linux

WORKDIR /src
COPY ./src /src
RUN go build -o /build/app  .

FROM alpine:latest
COPY ./env /env
COPY --from=builder /build/app .
ENTRYPOINT ["/app"]

