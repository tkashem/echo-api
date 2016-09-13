FROM scratch

ARG PIPELINE_LABEL="Unknown"
LABEL Pipeline=${PIPELINE_LABEL}

EXPOSE 3000

ADD artifacts/main /main
ENTRYPOINT ["/main"]
