FROM ubuntu:20.04
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update && apt-get install -y curl
RUN curl -L https://tarantool.io/gNXpbCs/release/2.6/installer.sh | bash
RUN apt-get install -y tarantool 

WORKDIR /tarantool

COPY scripts/lua/script.lua script.lua

EXPOSE 5555

ENTRYPOINT [ "tarantool", "script.lua" ]
