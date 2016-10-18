#!/bin/bash
set -e

function abort() {
  echo $@
  exit 1
}

function wait_for_port {
  port=$1
  if [[ -z $port ]]; then
    abort "USAGE: port"
  fi
  for i in $(seq 1 300); do
    if echo exit | nc localhost $port; then
      return
    fi
    sleep 0.1
  done
  abort "timeout waiting for port $port"
}

function wait_for_http {
  addr=$1
  if [[ -z $addr ]]; then
    abort "USAGE: address"
  fi
  for i in $(seq 1 300); do
    if curl -o /dev/null -sIf $addr; then
      return
    fi
    sleep 0.1
  done
  abort "timeout waiting for address $addr"
}

function install_dependencies {
  apk update
  apk add alpine-sdk jq curl postgresql git go
}

function setup_project {
  export TMP_PROJ=$(mktemp -d)
  echo created dir $TMP_PROJ

  export GOPATH=$TMP_PROJ
  export PROJ=$TMP_PROJ/src/github.com/tobstarr/k8s-cd-test

  mkdir -p $PROJ

  git clone https://github.com/tobstarr/k8s-cd-test.git $PROJ
  cd $PROJ
}

function wait_for_services {
  echo waiting for postgres
  wait_for_port 5432
  echo postgres available

  psql -h 127.0.0.1 -U postgres -c "CREATE DATABASE k8s_cd_test ENCODING 'utf8'"

  echo created database

  echo waiting for elasticsearch
  wait_for_http http://127.0.0.1:9200
  echo elasticsearch available
}

echo "running tests"

install_dependencies
setup_project
wait_for_services

packages=$(go list ./... | grep -v "/vendor")
go test -v $packages
