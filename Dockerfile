FROM alpine:3.23

RUN apk --no-cache add bash ca-certificates git openssh-client curl rsync docker-cli jq yq

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/skpr $TARGETPLATFORM/skpr-agent $TARGETPLATFORM/skpr-rsh /usr/local/bin/

CMD ["skpr"]
