#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")

set -e

cd "$script_path"

source ../../.env

cluster=ak0-cluster

declare -A hostmap

hostmap["${WORKER_NODE_0_HOSTNAME}"]=worker-node-0
hostmap["${WORKER_NODE_1_HOSTNAME}"]=worker-node-1

for host in "${!hostmap[@]}"; do
  echo $host
  echo "${hostmap[$host]}"
  kubectl config set-cluster "$cluster" \
    --certificate-authority=certs/ca.crt \
    --embed-certs=true \
    --server="https://${MASTER_NODE_0_HOSTNAME}.kubernetes.local:6443" \
    --kubeconfig="kubeconfigs/${host}.kubeconfig"

  kubectl config set-credentials "system:node:${host}" \
    --client-certificate="certs/${hostmap[$host]}.crt" \
    --client-key="certs/${hostmap[$host]}.key" \
    --embed-certs=true \
    --kubeconfig="kubeconfigs/${host}.kubeconfig"

  kubectl config set-context default \
    --cluster="$cluster" \
    --user="system:node:${host}" \
    --kubeconfig="kubeconfigs/${host}.kubeconfig"

  kubectl config use-context default \
    --kubeconfig="kubeconfigs/${host}.kubeconfig"
done

