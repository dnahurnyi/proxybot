############################
# STEP 1 install tdlib
############################

FROM alpine:3.12 as tdlib

WORKDIR /

RUN apk add --no-cache \
        ca-certificates

RUN apk add --no-cache --virtual .build-deps \
        g++ \
        make \
        cmake \
        git \
        gperf \
        libressl-dev \
        zlib-dev \
        zlib-static \
        linux-headers;

RUN git clone https://github.com/tdlib/td.git && \
    cd td && \
    git checkout v1.8.0 && \
    mkdir build && \
    cd build && \
    cmake -DCMAKE_BUILD_TYPE=Release .. && \
    cmake --build . && \
    make install

############################
# STEP 2 build go binary
############################

FROM golang:1.17-alpine AS builder
COPY --from=tdlib /usr/local/include/td /usr/local/include/td
COPY --from=tdlib /usr/local/lib/libtd* /usr/local/lib/
COPY --from=tdlib /usr/lib/libssl.a /usr/local/lib/libssl.a
COPY --from=tdlib /usr/lib/libcrypto.a /usr/local/lib/libcrypto.a
COPY --from=tdlib /lib/libz.a /usr/local/lib/libz.a
RUN apk add build-base
WORKDIR /app
COPY . ./
RUN go build --ldflags "-extldflags '-static -L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz'" -a -o ./artifacts/svc

############################
# STEP 3 build a small image
############################
FROM scratch

# Copy our static executable
COPY --from=builder /app/artifacts/svc /svc
COPY --from=builder /app/.tdlib /.tdlib

# Run the svc binary.
CMD ["./svc"]
