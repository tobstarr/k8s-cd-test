#!/bin/bash
set -e

function abort() {
	echo $@
	exit 1
}

function wait_for_http {
	addr=$1
	if [[ -z $addr ]]; then
		abort "USAGE: address"
	fi
	for i in $(seq 1 3000); do
		if curl -o /dev/null -sIf $addr; then
			return
		fi
		sleep 0.1
	done
	abort "timeout waiting for address $addr"
}

apk update
apk add curl

echo "waiting for app"
wait_for_http http://127.0.0.1:3000
echo "app available"

echo "waiting for selenium"
wait_for_http http://localhost:4444/wd/hub
echo "selenium available"

RUN_INTEGRATION_TESTS=true /usr/bin/k8s-cd-test
