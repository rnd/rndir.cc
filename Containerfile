FROM scratch
COPY ./assets /assets
COPY ./dist/site /
ENTRYPOINT ["/site"]
