FROM ubuntu:22.04

RUN apt-get update && apt-get -y upgrade && apt-get install -y ca-certificates && update-ca-certificates
RUN apt-get install -y tor obfs4proxy

COPY ./torrc /etc/tor/torrc

EXPOSE 9050

RUN groupadd tor && useradd -g tor tor
USER tor
WORKDIR /home/tor

ENTRYPOINT ["tor"]