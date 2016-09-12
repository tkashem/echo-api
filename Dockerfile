FROM ubuntu:14.04

ARG PIPELINE_LABEL="Unknown"
LABEL Pipeline=${PIPELINE_LABEL}

EXPOSE 3000
ADD artifacts/main /main
ENTRYPOINT ["/main"]
