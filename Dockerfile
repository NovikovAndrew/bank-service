FROM ubuntu:latest
LABEL authors="andreynovikov"

ENTRYPOINT ["top", "-b"]