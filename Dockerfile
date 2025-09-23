FROM alpine:3.21

RUN apk --no-cache add bash ca-certificates git openssh-client curl rsync docker-cli jq yq
COPY skpr skpr-rsh /usr/local/bin/

CMD ["skpr"]
