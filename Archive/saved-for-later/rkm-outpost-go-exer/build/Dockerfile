# build
FROM golang:alpine as builder
RUN apk update && apk add --no-cache
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    USER=appuser \
    UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/noexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY . .
RUN go build -ldflags="-w -s" -o /rkm-outpost

# final image
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /rkm-outpost /rkm-outpost

USER 10001
ENTRYPOINT ["./rkm-outpost"]