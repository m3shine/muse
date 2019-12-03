
#FROM golang:alpine AS build-env
#ENV PACKAGES make gcc bash
#WORKDIR /go/src/muse
#COPY . .
#RUN mkdir -p $HOME/go/bin && \
#    echo "export GOPATH=$HOME/go" >> ~/.bash_profile && \
#    echo "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile && \
#    echo "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile && \
#    echo "export GO111MODULE=on" >> ~/.bash_profile && \
#    source ~/.bash_profile && \
#    apk add --no-cache $PACKAGES && \
#    make install

FROM alpine:edge
WORKDIR /
COPY ./build/mused /usr/bin/mused
COPY ./build/musecli /usr/bin/musecli
#SHELL ["#!/bin/sh"]
CMD ["mused","start"]

# sed -i '.bak' "s/persistent_peers = .*/persistent_peers = 'id@first_node_ip:26656'/g" ~/.mused/config/config.toml
# run "musecli status" on first node to get id.