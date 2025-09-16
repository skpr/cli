ARG COMPILE_IMAGE=scratch
FROM ${COMPILE_IMAGE} as build

FROM skpr/nginx:1.x
ADD .skpr/package/nginx/custom.conf /etc/nginx/conf.d/custom.conf
ADD .skpr/package/nginx/redirects.conf /etc/nginx/redirects.conf
COPY --chown=skpr:skpr --from=build /data /data
RUN rm -rf /data/node_modules
