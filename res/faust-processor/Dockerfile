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
        librocksdb-dev \
        build-essential autoconf libtool pkg-config libpython3.8-dev \
        librocksdb-dev \
        libsnappy-dev \
        zlib1g-dev \
        libbz2-dev \
        liblz4-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && :

RUN python3.8 -m venv /venv
ENV PATH=/venv/bin:$PATH

RUN python -V
RUN python -m pip install faust[rocksdb] nltk
RUN python -c "import nltk; nltk.download('vader_lexicon')"

RUN mkdir /app
COPY util /app/util
COPY process /app/process

WORKDIR /app

ENTRYPOINT python -m process worker -l info