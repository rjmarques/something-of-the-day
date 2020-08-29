## Backend build
FROM something-backend-build-img as backend-build
RUN GOOS=linux GOARCH=amd64 go build

## Frontend build
FROM something-frontend-build-img as frontend-build
RUN mkdir build
RUN chmod -R 777 .
RUN ls -la
RUN whoami
RUN cp /root/workspace/something-of-the-day/public/favicon.ico /root/workspace/something-of-the-day/build/favicon.ico
RUN npm run build

## Final container that holds the artifacts
FROM alpine:3.11 
RUN apk update \
    && apk add --no-cache bash

WORKDIR /home/something-of-the-day
RUN chmod 777 -R .

COPY --from=backend-build /root/workspace/something-of-the-day/something-of-the-day .

COPY --from=frontend-build /root/workspace/something-of-the-day/build ./frontend

EXPOSE 80

ENTRYPOINT ["./something-of-the-day"]