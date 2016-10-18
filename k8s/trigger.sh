#!/bin/bash
set -e

function abort() {
  echo $@
  exit 1
}

file=$1

if [[ -z $file ]]; then
	echo "Usage: file"
	exit 1
fi

id=$(bash $file | kubectl create -f - | cut -d '"' -f 2)
echo started container $id

while true; do
  phase=$(kubectl get pods/${id} -o json | jq '.status.phase' -c -r)
  if [[ -z $phase || $phase == "Pending" ]]; then
    sleep 0.1
    continue
  fi
  break
done

while true; do
	if kubectl logs -f pod/${id} -c default 2> /dev/null; then
		break
	fi
	sleep 0.1
done
