#!/bin/sh
set -e
if [ -n "$1" ]; then
  if [ "$1" = "--build" ]; then
    docker build . -t hugo
    exit
  fi
fi
docker run -v$(pwd):/work --workdir=/work --entrypoint /bin/hugo --rm hugo $@
