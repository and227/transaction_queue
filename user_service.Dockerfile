FROM golang:1.19

EXPOSE 8080

RUN mkdir /app
WORKDIR /app

# RUN apt update && apt upgrade -y && \
#     apt install -y git \
#     make openssh-client

# RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
#     && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

COPY ./wait-for-it.sh ./wait-for-it.sh

CMD \
    ./wait-for-it.sh \
    ${POSTGRES_HOST:-localhost}:${POSTGRES_PORT:-5432} \
    ${REDIS_HOST:-localhost}:${REDIS_PORT:-6379} \
    --timeout=25 -- \
    ls -la ./user_service && \
    ./build/user_service_main
    # air