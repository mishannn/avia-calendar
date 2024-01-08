FROM ubuntu:22.04

RUN apt-get update && apt-get -y upgrade && apt-get install -y ca-certificates && update-ca-certificates
RUN apt-get install -y tor

COPY ./torrc /etc/tor/torrc

EXPOSE 9050

ENTRYPOINT ["tor"]