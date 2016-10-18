#!/bin/bash
set -e

cmd="docker build"

if [[ -n $IMAGE_TAG ]]; then
  cmd="$cmd -t $IMAGE_TAG"
fi

echo $cmd
$cmd .

if [[ -n $IMAGE_TAG ]]; then
  docker push $IMAGE_TAG
fi
