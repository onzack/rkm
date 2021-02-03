# build
FROM golang:alpine as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY . .
RUN go build -o /rkm-outpost

# final image
FROM alpine:latest
RUN addgroup -S speedtest && adduser -S speedtest -G speedtest
COPY --from=builder /rkm-outpost /rkm-outpost
RUN chmod +x /rkm-outpost
USER speedtest
ENTRYPOINT ["/rkm-outpost"]