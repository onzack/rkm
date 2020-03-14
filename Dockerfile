FROM ubuntu:19.10
LABEL maintainer="dmlabs.ch"
RUN groupadd -g 10000 sentry-outpost && useradd -u 10000 -g 10000 -s /bin/bash -m sentry-outpost
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends tzdata curl ca-certificates jq gnupg gnupg2 gnupg1 apt-transport-https bc && \
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends kubectl  && \
    apt-get purge -y gnupg gnupg2 gnupg1 apt-transport-https && \
    apt-get -y autoremove && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
COPY --chown=sentry-outpost:sentry-outpost collect-metrics.bash /opt/collect-metrics.bash
USER sentry-outpost
WORKDIR /home/sentry-outpost
ENTRYPOINT ["/opt/collect-metrics.bash"]