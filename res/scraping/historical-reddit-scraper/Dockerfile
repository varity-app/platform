# Fetch till from github
FROM debian:buster-slim as prod

RUN : \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends wget ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && :

WORKDIR /opt
RUN wget https://github.com/DataHenHQ/till/releases/download/v0.8.0/till_0.8.0_Linux_x86_64.tar.gz
RUN tar -xf till_0.8.0_Linux_x86_64.tar.gz

ENTRYPOINT /opt/till serve -t ${TILL_TOKEN}