#!/bin/sh
set -eux
IMAGE_NAME=synmodoro_app
docker build . -t $IMAGE_NAME
# docker run --entrypoint='' $IMAGE_NAME /goapp/app
