#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")

set -e

cd "${script_path}"

source ../../.env

echo ""
echo "running install.sh"

download_dir="/root/downloads"
bin_dir="/usr/local/bin"
kubernetes_dir="/etc/kubernetes"
kubernetes_pki_dir="/etc/kubernetes/pki"

nodes=(
    "$MASTER_NODE_0_HOSTNAME"
    # "$WORKER_NODE_0_HOSTNAME" 
    # "$WORKER_NODE_1_HOSTNAME"
    # "$LOAD_BALANCER_HOSTNAME"
)

for host in "${nodes[@]}"; do
  # set up .bashrc
  scp bashrc "$host:/root/.bashrc"

  if [[ "$host" != "$LOAD_BALANCER_HOSTNAME" ]];then
    ssh "$host" mkdir -p $download_dir $kubernetes_dir $kubernetes_pki_dir

    scp ../../configs/kubernetes/pki/kubeconfigs/kube-proxy.conf \
      "${host}:${kubernetes_dir}/kube-proxy.conf"

    # disable swap
    ssh "$host" << 'EOF'
    set -e
    apt-get -y install jq
    swapoff -a
    sed -i 's/^.*swap.*$/# &/' /etc/fstab
    for unit in $(systemctl --type swap --all --output=json | jq -r '.[] | .unit'); do
      systemctl mask $unit;
    done
EOF

    if [[ $host = "$MASTER_NODE_0_HOSTNAME" ]]; then
      scp ../../configs/kubernetes/pki/certs/ca.crt \
        ../../configs/kubernetes/pki/certs/ca.key \
        ../../configs/kubernetes/pki/kubeconfigs/kube-controller-manager.conf \
        ../../configs/kubernetes/pki/kubeconfigs/kube-proxy.conf \
        ../../configs/kubernetes/pki/kubeconfigs/kube-scheduler.conf \
        ../../configs/kubernetes/encryption-config.yaml "${host}:${download_dir}/"

      ssh "$host" << EOF
      set -e
      mv "${download_dir}/ca.crt" "${download_dir}/ca.key" "$kubernetes_pki_dir"
      mv "${download_dir}/kube-proxy.conf" \
        "${download_dir}/kube-controller-manager.conf" \
        "${download_dir}/kube-scheduler.conf" "$kubernetes_dir"

      export ENCRYPTION_KEY="$(head -c 32 /dev/urandom | base64)"
      envsubst < "${download_dir}"/encryption-config.yaml \
        > "${kubernetes_dir}/encryption-config.yaml"
EOF
    fi

    # copy up crio, crun, and conmon
    crun_bin=crun-1.23.1-linux-amd64
    conmon_bin=conmon.amd64
    scp -r ../vendor/crio/ \
      "../vendor/crun/${crun_bin}" \
      "../vendor/conmon/${conmon_bin}" \
      "${host}:${download_dir}/"

    # mv things into place, install crio`
    ssh "$host" << EOF
    set -e
    mv "${download_dir}/${crun_bin}" \
      "${download_dir}/${conmon_bin}" "${bin_dir}"
    systemctl enable --now nftables.service
    # Load the bridge netfilter module
    modprobe br_netfilter
    echo "br_netfilter" >> /etc/modules-load.d/modules.conf

    # configure iptables. note debian translates to nftables under the hood
    echo "net.bridge.bridge-nf-call-iptables = 1" >> /etc/sysctl.d/kubernetes.conf
    echo "net.bridge.bridge-nf-call-ip6tables = 1" >> /etc/sysctl.d/kubernetes.conf
    echo "net.ipv4.ip_forward = 1"  >> /etc/sysctl.d/kubernetes.conf
    sysctl -p /etc/sysctl.d/kubernetes.conf

    pushd ${download_dir}/crio >/dev/null 2>&1
    ./install
    systemctl enable --now crio.service
EOF

    # install, configure, and enable kubelet
    scp ../vendor/kubernetes/kubelet \
      ../../configs/kubernetes/kubelet.service \
      ../../configs/kubernetes/kubeapi-to-kubelet.yaml \
      ../../configs/kubernetes/kubelet-config.yaml \
      "../../configs/kubernetes/pki/kubeconfigs/${host}.conf" "${host}:${download_dir}"

    ssh "$host" << EOF
    set -e
    mv "${download_dir}/kubelet" "${bin_dir}/"
    mv "${download_dir}/kubelet.service" /etc/systemd/system/
    mv "${download_dir}/kubelet-config.yaml" /tmp/
    mv "${download_dir}/kubeapi-to-kubelet.yaml" "${kubernetes_dir}"
    mv "${download_dir}/${host}.conf" "${kubernetes_dir}/kubelet.conf"
EOF
    ssh "$host" << 'EOF'
    set -e
    export KUBELET_INTERNAL_NODE_IP=$(ip addr show eth1 | grep -oP 'inet \K[0-9.]+' | head -1)
    envsubst < /tmp/kubelet-config.yaml > /etc/kubernetes/kubelet-config.yaml
    systemctl daemon-reload
    systemctl enable --now kubelet.service
EOF
  fi
done

