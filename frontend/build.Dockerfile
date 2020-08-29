FROM node:alpine3.11
RUN apk update \
    && apk add --no-cache bash

ADD . /root/workspace/something-of-the-day
WORKDIR /root/workspace/something-of-the-day

# for circle ci builds
RUN chmod 777 -R .

RUN npm install