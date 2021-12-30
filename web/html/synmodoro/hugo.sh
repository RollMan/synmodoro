#!/bin/sh
docker run -v$(pwd):/work --workdir=/work --entrypoint /bin/hugo --rm hugo $@
