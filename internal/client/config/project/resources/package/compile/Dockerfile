FROM skpr/node:14-1.x as frontend
ADD --chown=skpr:skpr . /data
RUN make npm-ci styleguide

FROM skpr/php-cli:8.0-1.x as backend
ADD --chown=skpr:skpr . /data
RUN make composer-prod

FROM scratch
COPY --from=frontend /data /data
COPY --from=backend /data /data
