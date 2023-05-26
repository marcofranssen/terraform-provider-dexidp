#!/usr/bin/env bash

set -e

K8S_NS="${1:-terraform-provider-dexidp}"

trap 'trap - SIGTERM && kill -- $PID' SIGINT SIGTERM EXIT

kubectl -n "${K8S_NS}" port-forward services/dex 5557:grpc >/dev/null &
PID=$!

TF_ACC=1 go test -v -count=1 ./pkg/dexidp
