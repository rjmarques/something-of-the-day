FROM node:12.18.3-alpine3.11
RUN apk update \
    && apk add --no-cache bash

ADD . /root/workspace/something-of-the-day
WORKDIR /root/workspace/something-of-the-day

RUN npm install

RUN mkdir build
RUN chmod -R 777 .
RUN ls -la .
RUN npm cache clean --force
RUN npm run build