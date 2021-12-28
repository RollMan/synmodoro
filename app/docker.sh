#!/bin/sh
set -eux
IMAGE_NAME=synmodoro_app
IMAGE_NAME=golang:1.14.7-alpine3.12
docker run --workdir=/work -v$(pwd):/work -it $IMAGE_NAME $@
