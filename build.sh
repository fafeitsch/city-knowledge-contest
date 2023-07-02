#!/bin/sh
releaseversion=$(git tag --sort=-version:refname | head -n 1)
releasetime=$(date +"%Y%m%d%H%M")
docker build --rm --no-cache --build-arg VERSION=$releaseversion+$releasetime -t  city-contest .
