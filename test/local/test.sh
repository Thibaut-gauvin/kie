#!/bin/bash

set -eu

CONTEXT="kind-kind"
NAMESPACE="default"
DIRNAME=$(dirname "$0")

function apply() {
  echo "applying $1"
  kubectl --context="$CONTEXT" --namespace="$NAMESPACE" apply -f $1
}

function cp() {
  local local_file=$1
  local pod=$2
  local remote_file=$3
  echo "cp $local_file into $pod:$remote_file"
  kubectl --context="$CONTEXT" --namespace="$NAMESPACE" cp "$local_file" "$pod:$remote_file"
}

function exec() {
  local pod=$1
  local exec_cmd=$2
  echo "exec into $pod run $exec_cmd"
  kubectl --context="$CONTEXT" --namespace="$NAMESPACE" exec -it "$pod" -- bash -c "$exec_cmd"
}

# Create service-account & cluster-role & cluster-role-binding
#apply "$DIRNAME/debug_resources.yml"

# Check service-account authorisations
#kubectl auth can-i --list --as=system:serviceaccount:default:kie

#apply "$DIRNAME/debug_pod.yml"
#cp "$DIRNAME/../../dist/kie_linux_amd64_v1/kie" "debugpod" "/tmp/kie"
#exec "debugpod" "./tmp/kie serve -l debug"
#kubectl delete pod debugpod --grace-period=5
