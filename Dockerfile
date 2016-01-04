FROM scratch

ADD build/dispatch /
ADD ca-certificates.crt /etc/ssl/certs/

VOLUME ["/data"]

ENTRYPOINT ["/dispatch"]
CMD ["--dir=/data"]
