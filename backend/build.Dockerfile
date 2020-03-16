FROM alpine:3.11

# Install the major dependencies
RUN apk update \
    && apk add --no-cache bash build-base util-linux git openssh go

# copy the local files to the container's workspace
ADD . /root/workspace/something-of-the-day
WORKDIR /root/workspace/something-of-the-day
