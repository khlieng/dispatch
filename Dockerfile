FROM scratch

ADD build/name_pending /
ADD ca-certificates.crt /etc/ssl/certs/

VOLUME ["/data"]

ENTRYPOINT ["/name_pending"]
CMD ["-p=8080", "--dir=/data"]