#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")

set -e

source .env

cd "$script_path"


echo ""
echo "running export_certs.sh"


# copy certs to worker nodes
for host in $WORKER_NODE_0_HOSTNAME $WORKER_NODE_1_HOSTNAME; do
  ssh root@"${host}" mkdir -p /var/lib/kubelet/

  scp certs/ca.crt root@"${host}":/var/lib/kubelet/

  scp "certs/${host}.crt" \
    root@"${host}":/var/lib/kubelet/kubelet.crt

  scp "certs/${host}.key" \
    root@"${host}":/var/lib/kubelet/kubelet.key
done

# copy certs to api-server
scp \
  ca.key ca.crt \
  kube-api-server.key kube-api-server.crt \
  service-accounts.key service-accounts.crt \
  root@"$MASTER_NODE_0_HOSTNAME":~/
