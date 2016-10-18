#!/bin/bash
set -e

apk update && apk add git go

repo=github.com/tobstarr/k8s-cd-test
dir=$(mktemp -d)
project=$dir/src/${repo}

mkdir -p $project

git clone https://${repo} $project

cd $project

rev=$(git rev-list HEAD -n 1 | cut -b 1-8)

export IMAGE_TAG=127.0.0.1:30015/tobstarr/k8s-cd-test:$rev

echo using IMAGE_TAG=$IMAGE_TAG

sh ./docker_build.sh

echo "pushed image $IMAGE_TAG"
