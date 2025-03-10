FROM alpine:3

RUN apk add --update curl ca-certificates && rm -rf /var/cache/apk* # Certificates for SSL

COPY bin/bedrock-sidecar .
ENTRYPOINT [ "./bin/bedrock-sidecar" ]