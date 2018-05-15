FROM centos:latest

ARG PIPELINE_LABEL="Unknown"
LABEL Pipeline=${PIPELINE_LABEL}

EXPOSE 3000

ADD artifacts/echo /echo
ENTRYPOINT ["/echo"]
