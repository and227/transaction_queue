FROM golang:1.19

EXPOSE 8080

RUN mkdir /app
WORKDIR /app

COPY ./wait-for-it.sh ./wait-for-it.sh

CMD \
    ./wait-for-it.sh \
    ${POSTGRES_HOST:-localhost}:${POSTGRES_PORT:-5432} \
    ${REDIS_HOST:-localhost}:${REDIS_PORT:-6379} \
    --timeout=25 -- \
    ls -la ./user_service && \
    ./build/user_service_main
