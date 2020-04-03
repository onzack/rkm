FROM alpine:latest
LABEL maintainer="https://github.com/dmlabs"
RUN addgroup -g 10000 -S patch-operator && adduser -u 10000 -S -G patch-operator patch-operator
WORKDIR /home/patch-operator
RUN apk update && \
    apk upgrade && \
    apk add tzdata curl bash && \
    curl -o /usr/local/bin/kubectl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod -R +x /usr/local/bin && \
    chown -R patch-operator:patch-operator /usr/local/bin && \
    apk del --purge curl && \
    rm -f /var/cache/apk/*
COPY --chown=patch-operator:patch-operator patch-operator.bash /usr/local/bin/patch-operator.bash
USER patch-operator
ENTRYPOINT ["/usr/local/bin/patch-operator.bash"]