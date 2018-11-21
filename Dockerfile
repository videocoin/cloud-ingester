FROM bitnami/minideb:jessie as builder

LABEL maintainer="Dmitry Serkin <dmitry@liveplanet.net>"

ENV NGINX_VERSION 1.13.6
ENV NGINX_URL https://github.com/nginx/nginx/archive
ENV NGINX_RTMP_MODULE_VERSION 1.1.7.16
ENV NGINX_RTMP_MODULE_URL https://github.com/reality-lab-networks/nginx-rtmp-module/archive

# Install dependencies
RUN apt-get update && apt-get install -y \
    dos2unix \
    mercurial \
    libreadline-dev \
    libncurses5-dev \
    libpcre3-dev \
    libssl-dev \
    perl \
    make \
    musl \
    build-essential \
    wget \
    ca-certificates \
    openssl \
    bzip2 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /opt/

ADD . ./

# Download and decompress Nginx
RUN mkdir -p /tmp/build/nginx
RUN wget --quiet -O /tmp/build/nginx-${NGINX_VERSION}.tar.gz ${NGINX_URL}/release-${NGINX_VERSION}.tar.gz
RUN tar -C /tmp/build/nginx -zxvf /tmp/build/nginx-${NGINX_VERSION}.tar.gz
RUN rm /tmp/build/nginx-${NGINX_VERSION}.tar.gz
RUN ls /tmp/build/nginx

# Download and decompress RTMP module
RUN mkdir -p /tmp/build/nginx-rtmp-module
RUN wget --quiet -O /tmp/build/nginx-rtmp-module/nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION}.tar.gz ${NGINX_RTMP_MODULE_URL}/v${NGINX_RTMP_MODULE_VERSION}.tar.gz
RUN tar -C /tmp/build/nginx-rtmp-module -zxf /tmp/build/nginx-rtmp-module/nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION}.tar.gz
RUN rm /tmp/build/nginx-rtmp-module/nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION}.tar.gz

# Download FFMPEG binary
RUN cd /tmp
RUN wget --quiet -O /tmp/ffmpeg-04dc3b4.tbz2 https://storage.googleapis.com/liveplanet-releases/ffmpeg/master/ffmpeg-04dc3b4.tbz2
RUN tar xjf /tmp/ffmpeg-04dc3b4.tbz2
RUN cp ffmpeg /usr/bin/ffmpeg
# RUN cp -r lib/* /lib

# Build nginx
RUN cd /tmp/build/nginx/nginx-release-${NGINX_VERSION} && \
    auto/configure --prefix=/opt/stream-ingester --with-debug --with-http_ssl_module --add-module=/tmp/build/nginx-rtmp-module/nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION} && \
    make && \
    make install

FROM bitnami/minideb:jessie AS release

COPY --from=builder /opt /opt
COPY --from=builder /opt/etc /opt/stream-ingester/etc
COPY --from=builder /opt/src /opt/stream-ingester/src
COPY --from=builder /opt/var /opt/stream-ingester/var
COPY --from=builder /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=builder /lib /lib

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
    libtbb-dev \
    curl \
    musl \
    libssl1.0.0 && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /tmp/records /var/log/stream-ingester
RUN chown www-data /tmp/records /var/log/stream-ingester

EXPOSE 80
EXPOSE 1935
EXPOSE 8888

RUN mkdir -p /var/log/stream-ingester

ENV PATH="/opt/stream-ingester/sbin:${PATH}"

ENTRYPOINT ["nginx"]
CMD ["-c", "/opt/stream-ingester/etc/stream-ingester-rtmp.conf", "-g", "daemon off;"]