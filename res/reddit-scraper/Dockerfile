FROM ubuntu:18.04

RUN apt update -y
RUN apt install -y software-properties-common

RUN : \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        software-properties-common \
    && add-apt-repository -y ppa:deadsnakes \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        python3.8-venv \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && :

RUN python3.8 -m venv /venv
ENV PATH=/venv/bin:$PATH

RUN python -V
RUN python -m pip install asyncpraw google-cloud-pubsub google-cloud-firestore

RUN mkdir /app
COPY entrypoints/scrape_reddit_live.py /app
COPY util /app/util
COPY scrapers /app/scrapers
COPY messages /app/messages

WORKDIR /app

ENTRYPOINT python scrape_reddit_live.py $MODE --sleep $SLEEP --subreddits $SUBREDDITS