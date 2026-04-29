FROM gcr.io/distroless/static-debian12:nonroot
COPY selectronic_exporter /bin/selectronic_exporter
EXPOSE 9788
USER nonroot:nonroot
ENTRYPOINT ["/bin/selectronic_exporter"]
