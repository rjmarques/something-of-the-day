## Backend build
FROM something-backend-build-img as backend-build
RUN GOOS=linux GOARCH=amd64 go build

## Frontend build
FROM something-frontend-build-img as frontend-build
RUN npm run build

## Final container that holds the artifacts
FROM alpine:3.11 
RUN apk update \
    && apk add --no-cache bash

WORKDIR /home/something-of-the-day

COPY --from=backend-build /root/workspace/something-of-the-day/something-of-the-day .

COPY --from=frontend-build /root/workspace/something-of-the-day/build ./frontend/build

EXPOSE 80

ENTRYPOINT ["./something-of-the-day"]