FROM golang:latest AS builder

RUN apt update && apt upgrade -y

WORKDIR /opt/

ARG SSH_PRIVATE_KEY
RUN mkdir /root/.ssh/
RUN echo ${SSH_PRIVATE_KEY} > /root/.ssh/id_rsa
RUN echo "    IdentityFile /root/.ssh/id_rsa" >> /etc/ssh/ssh_config
RUN echo " Host gitlab.videocoin.io" >> /root/.ssh/config
RUN touch /root/.ssh/known_hosts
RUN chmod 600 /root/.ssh/id_rsa
RUN ssh-keyscan gitlab.videocoin.io >> /root/.ssh/known_hosts

ADD . ./


RUN make build

FROM ubuntu:latest AS release

COPY --from=builder /opt/bin/hookd ./

ENTRYPOINT ./hookd

