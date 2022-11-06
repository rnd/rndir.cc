FROM ubuntu:22.10
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update -q && apt-get install -qy \
    curl jq \
    texlive-full \
    python3-pygments gnuplot \
    make git \
    fonts-inconsolata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /data
VOLUME ["/data"]
