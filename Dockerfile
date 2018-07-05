FROM golang:1.10-stretch
MAINTAINER Nick Lange <nick.lange@gmail.com>



ENV SRC_DIR /go/src/github.com/NickJLange/sysInfoWebServer
ENV GOOS linux
ENV GOARCH amd64

COPY . $SRC_DIR
WORKDIR $SRC_DIR
RUN make linux TOPLEVEL=/go/

ENV SUEXEC_VERSION v0.2
ENV TINI_VERSION v0.16.1
RUN set -x \
  && cd /tmp \
  && git clone https://github.com/ncopa/su-exec.git \
  && cd su-exec \
  && git checkout -q $SUEXEC_VERSION \
  && make \
  && cd /tmp \
  && wget -q -O tini https://github.com/krallin/tini/releases/download/$TINI_VERSION/tini \
  && chmod +x tini


# Get the TLS CA certificates, they're not provided
RUN apt-get update && apt-get install -y ca-certificates

FROM busybox:1-glibc
MAINTAINER Nick Lange <nick.lange@gmail.com>

ENV SRC_DIR /go/src/github.com/NickJLange/sysInfoWebServer

COPY --from=0 /go//pkg//linux_arm/github.com/NickJLange/sysInfoWebServer//sysInfoWebServer /usr/local/bin/

COPY --from=0 /tmp/su-exec/su-exec /sbin/su-exec
COPY --from=0 /tmp/tini /sbin/tini
COPY --from=0 /etc/ssl/certs /etc/ssl/certs

# This shared lib (part of glibc) doesn't seem to be included with busybox.
COPY --from=0 /lib/x86_64-linux-gnu/libdl-2.24.so /lib/libdl.so.2

EXPOSE 8008

ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/sysInfoWebServer"]

# Execute the daemon subcommand by default
CMD ["daemon", "--migrate=true"]
