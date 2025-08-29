ARG COMPILE_IMAGE=scratch
FROM ${COMPILE_IMAGE} as build

FROM skpr/php-fpm:8.0-1.x
ADD .skpr/package/fpm/custom.ini /etc/php/conf.d/50_custom.ini
COPY --chown=skpr:skpr --from=build /data /data
RUN rm -rf /data/node_modules

USER root
RUN curl -sSL https://github.com/skpr/mail/releases/download/v0.0.7/skprmail_linux_amd64 -o /usr/local/bin/skprmail && \
    chmod +rx /usr/local/bin/skprmail
USER skpr
