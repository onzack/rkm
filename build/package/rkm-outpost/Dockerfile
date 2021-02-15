# build
FROM golang:alpine as builder
ENV USER=rkm-outpost \
    UID=10000

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/noexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

COPY ./rkm-outpost /go/bin/rkm-outpost

# final image
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/rkm-outpost /go/bin/rkm-outpost

USER 10000:10000
ENTRYPOINT ["/go/bin/rkm-outpost"]