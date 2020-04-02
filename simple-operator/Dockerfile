FROM alpine:latest
LABEL maintainer="https://github.com/dmlabs"
RUN addgroup -g 10000 -S simple-operator && adduser -u 10000 -S -G simple-operator simple-operator
WORKDIR /home/simple-operator
RUN apk update && \
    apk upgrade && \
    apk add tzdata curl bash && \
    curl -o /usr/local/bin/kubectl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod -R +x /usr/local/bin && \
    chown -R simple-operator:simple-operator /usr/local/bin && \
    apk del --purge curl && \
    rm -f /var/cache/apk/*
COPY --chown=simple-operator:simple-operator simple-operator.bash /usr/local/bin/simple-operator.bash
USER simple-operator
ENTRYPOINT ["/usr/local/bin/simple-operator.bash"]