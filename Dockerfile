############################
# STEP 1 build go binary
############################

FROM golang:1.18-alpine AS builder
COPY --from=n826/tdlib:v1 /usr/local/include/td /usr/local/include/td
COPY --from=n826/tdlib:v1 /usr/local/lib/libtd* /usr/local/lib/
COPY --from=n826/tdlib:v1 /usr/lib/libssl.a /usr/local/lib/libssl.a
COPY --from=n826/tdlib:v1 /usr/lib/libcrypto.a /usr/local/lib/libcrypto.a
COPY --from=n826/tdlib:v1 /lib/libz.a /usr/local/lib/libz.a
RUN apk add build-base
WORKDIR /app
COPY . ./
RUN go build --ldflags "-extldflags '-static -L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz'" -a -o ./artifacts/svc

############################
# STEP 2 build a small image
############################
FROM scratch

# Copy our static executable
COPY --from=builder /app/artifacts/svc /svc
COPY --from=builder /app/.tdlib /.tdlib

# Run the svc binary.
CMD ["./svc"]
