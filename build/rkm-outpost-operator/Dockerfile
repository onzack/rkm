FROM alpine:latest
LABEL maintainer="https://github.com/dmlabs"
RUN addgroup -g 10000 -S rkm-operator && adduser -u 10000 -S -G rkm-operator rkm-operator
WORKDIR /home/rkm-operator
RUN apk update && \
    apk upgrade && \
    apk add tzdata curl bash && \
    curl -o /usr/local/bin/kubectl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod -R +x /usr/local/bin && \
    chown -R rkm-operator:rkm-operator /usr/local/bin && \
    apk del --purge curl && \
    rm -f /var/cache/apk/*
COPY --chown=rkm-operator:rkm-operator rkm-outpost-operator.bash /usr/local/bin/rkm-outpost-operator.bash
USER rkm-operator
ENTRYPOINT ["/usr/local/bin/rkm-outpost-operator.bash"]