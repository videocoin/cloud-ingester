FROM bitnami/minideb:jessie as builder

LABEL maintainer="Dmitry Serkin <dmitry@liveplanet.net>"

ENV OPENRESTY_VERSION 1.13.6.2
ENV OPENRESTY_URL https://github.com/openresty/openresty/archive
ENV NGINX_RTMP_MODULE_VERSION 1.1.7.13
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

# Download and decompress Openresty
RUN mkdir -p /tmp/build/openresty && \
    cd /tmp/build/openresty && \
    wget --quiet -O openresty-${OPENRESTY_VERSION}.tar.gz ${OPENRESTY_URL}/v${OPENRESTY_VERSION}.tar.gz && \
    tar -zxf openresty-${OPENRESTY_VERSION}.tar.gz && \
    rm openresty-${OPENRESTY_VERSION}.tar.gz

# Download and decompress RTMP module
RUN mkdir -p /tmp/build/nginx-rtmp-module && \
    cd /tmp/build/nginx-rtmp-module && \
    wget --quiet -O nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION}.tar.gz ${NGINX_RTMP_MODULE_URL}/v${NGINX_RTMP_MODULE_VERSION}.tar.gz && \
    tar -zxf nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION}.tar.gz && \
    rm nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION}.tar.gz

ADD ./ /usr/src/stream-ingester

# Download FFMPEG binary
RUN cd /tmp && \
    wget --quiet https://storage.googleapis.com/liveplanet-releases/ffmpeg/master/ffmpeg-04dc3b4.tbz2 && \
    tar xjf /tmp/ffmpeg-04dc3b4.tbz2 && \
    cp ffmpeg /usr/bin/ffmpeg && \
    cp -r lib/* /lib

# Build OpenResty
RUN cd /tmp/build/openresty/openresty-${OPENRESTY_VERSION} && \
    make all && \
    cd ./openresty-${OPENRESTY_VERSION} && \
    ./configure \
	--prefix=/opt/stream-ingester \
	--with-debug \
	--with-luajit \
	--add-module=/tmp/build/nginx-rtmp-module/nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION} && \
    make && \
    make install

FROM bitnami/minideb:jessie

COPY --from=builder /opt /opt
COPY --from=builder /usr/src/stream-ingester/etc /opt/stream-ingester/etc
COPY --from=builder /usr/src/stream-ingester/src /opt/stream-ingester/src
COPY --from=builder /usr/src/stream-ingester/var /opt/stream-ingester/var
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

RUN mkdir -p /tmp/records
RUN chown www-data /tmp/records

EXPOSE 80
EXPOSE 1935

RUN mkdir -p /var/log/stream-ingester

ENV PATH="/opt/stream-ingester/nginx/sbin:${PATH}"

ENTRYPOINT ["nginx"]
CMD ["-c", "/opt/stream-ingester/etc/stream-ingester-rtmp.conf"]
