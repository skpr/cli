ARG COMPILE_IMAGE=scratch
FROM ${COMPILE_IMAGE} as build

FROM skpr/php-cli:8.0-1.x
ADD .skpr/package/cli/custom.ini /etc/php/conf.d/custom.ini
COPY --chown=skpr:skpr --from=build /data /data
RUN rm -rf /data/node_modules
