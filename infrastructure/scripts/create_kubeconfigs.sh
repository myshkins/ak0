#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")

set -e

cd "${script_path}/../../configs/kubernetes/pki"

source ../../../.env

rm -f kubeconfigs/*

cluster=ak0-cluster

declare -A hostmap

hostmap["${MASTER_NODE_0_HOSTNAME}"]=master-node-0
hostmap["${WORKER_NODE_0_HOSTNAME}"]=worker-node-0
hostmap["${WORKER_NODE_1_HOSTNAME}"]=worker-node-1

# worker nodes
for host in "${!hostmap[@]}"; do
  kubectl config set-cluster "$cluster" \
    --certificate-authority=certs/ca.crt \
    --embed-certs=true \
    --server="https://${MASTER_NODE_0_HOSTNAME}.kubernetes.local:6443" \
    --kubeconfig="kubeconfigs/${host}.conf"

  kubectl config set-credentials "system:node:${host}" \
    --client-certificate="certs/${hostmap[$host]}.crt" \
    --client-key="certs/${hostmap[$host]}.key" \
    --embed-certs=true \
    --kubeconfig="kubeconfigs/${host}.conf"

  kubectl config set-context default \
    --cluster="$cluster" \
    --user="system:node:${host}" \
    --kubeconfig="kubeconfigs/${host}.conf"

  kubectl config use-context default \
    --kubeconfig="kubeconfigs/${host}.conf"
done

# kube-proxy
kubectl config set-cluster "$cluster" \
  --certificate-authority=certs/ca.crt \
  --embed-certs=true \
  --server="https://${MASTER_NODE_0_HOSTNAME}.kubernetes.local:6443" \
  --kubeconfig=kubeconfigs/kube-proxy.conf

kubectl config set-credentials system:kube-proxy \
  --client-certificate=certs/kube-proxy.crt \
  --client-key=certs/kube-proxy.key \
  --embed-certs=true \
  --kubeconfig=kubeconfigs/kube-proxy.conf

kubectl config set-context default \
  --cluster="$cluster" \
  --user=system:kube-proxy \
  --kubeconfig=kubeconfigs/kube-proxy.conf

kubectl config use-context default \
  --kubeconfig=kubeconfigs/kube-proxy.conf

#kube-controller manager
kubectl config set-cluster "$cluster" \
  --certificate-authority=certs/ca.crt \
  --embed-certs=true \
  --server="https://${MASTER_NODE_0_HOSTNAME}.kubernetes.local:6443" \
  --kubeconfig=kubeconfigs/kube-controller-manager.conf

kubectl config set-credentials system:kube-controller-manager \
  --client-certificate=certs/kube-controller-manager.crt \
  --client-key=certs/kube-controller-manager.key \
  --embed-certs=true \
  --kubeconfig=kubeconfigs/kube-controller-manager.conf

kubectl config set-context default \
  --cluster="$cluster" \
  --user=system:kube-controller-manager \
  --kubeconfig=kubeconfigs/kube-controller-manager.conf

kubectl config use-context default \
  --kubeconfig=kubeconfigs/kube-controller-manager.conf

# kube-scheduler
kubectl config set-cluster "$cluster" \
  --certificate-authority=certs/ca.crt \
  --embed-certs=true \
  --server="https://${MASTER_NODE_0_HOSTNAME}.kubernetes.local:6443" \
  --kubeconfig=kubeconfigs/kube-scheduler.conf

kubectl config set-credentials system:kube-scheduler \
  --client-certificate=certs/kube-scheduler.crt \
  --client-key=certs/kube-scheduler.key \
  --embed-certs=true \
  --kubeconfig=kubeconfigs/kube-scheduler.conf

kubectl config set-context default \
  --cluster="$cluster" \
  --user=system:kube-scheduler \
  --kubeconfig=kubeconfigs/kube-scheduler.conf

kubectl config use-context default \
  --kubeconfig=kubeconfigs/kube-scheduler.conf

# admin user
kubectl config set-cluster "$cluster" \
  --certificate-authority=certs/ca.crt \
  --embed-certs=true \
  --server="https://${MASTER_NODE_0_HOSTNAME}.kubernetes.local:6443" \
  --kubeconfig=kubeconfigs/admin.conf

kubectl config set-credentials admin \
  --client-certificate=certs/admin.crt \
  --client-key=certs/admin.key \
  --embed-certs=true \
  --kubeconfig=kubeconfigs/admin.conf

kubectl config set-context default \
  --cluster="$cluster" \
  --user=admin \
  --kubeconfig=kubeconfigs/admin.conf

kubectl config use-context default \
  --kubeconfig=kubeconfigs/admin.conf
