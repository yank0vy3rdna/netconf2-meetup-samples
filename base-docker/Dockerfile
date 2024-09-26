FROM alpine:latest

RUN apk fix && \
    apk --no-cache --update add git cmake clang clang-dev make gcc g++ \
    libc-dev linux-headers pcre2-dev mbedtls-dev libssh-dev curl-dev openssl bash

RUN git clone -b v3.1.0 --depth 1 https://github.com/CESNET/libyang.git && \
    mkdir libyang/build && cd libyang/build && \
    cmake .. && \
    make install && \
    cd ../.. && \
    ldconfig /etc/ld.so.conf.d && \
    rm -rf libyang

RUN git clone -b v2.10.1 --depth 1 https://github.com/sysrepo/sysrepo.git && \
    mkdir sysrepo/build && cd sysrepo/build && \
    cmake .. && \
    make install && \
    cd ../.. && \
    ldconfig /etc/ld.so.conf.d && \
    rm -rf sysrepo

RUN git clone -b v3.3.3 --depth 1 https://github.com/CESNET/libnetconf2.git && \
    mkdir libnetconf2/build && cd libnetconf2/build && \
    cmake .. && \
    make install && \
    cd ../.. && \
    ldconfig /etc/ld.so.conf.d && \
    rm -rf libnetconf2

RUN git clone -b v2.2.28 --depth 1 https://github.com/CESNET/netopeer2.git && \
    mkdir netopeer2/build && cd netopeer2/build && \
    cmake .. && \
    make install && \
    cd ../.. && \
    rm -rf netopeer2
