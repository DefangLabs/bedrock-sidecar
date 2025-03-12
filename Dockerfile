FROM alpine:3

RUN apk add --update curl ca-certificates && rm -rf /var/cache/apk* # Certificates for SSL

COPY bin/bedrock-sidecar ./bin/
ENTRYPOINT [ "./bin/bedrock-sidecar" ]