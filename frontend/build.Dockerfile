FROM node:12.18.3-alpine3.11
RUN apk update \
    && apk add --no-cache bash

ADD . /root/workspace/something-of-the-day
WORKDIR /root/workspace/something-of-the-day

RUN npm install