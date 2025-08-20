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

nodes=(
    # "$MASTER_NODE_0_HOSTNAME"
    "$WORKER_NODE_0_HOSTNAME" 
    # "$WORKER_NODE_1_HOSTNAME"
    # "$LOAD_BALANCER_HOSTNAME"
)

for host in "${nodes[@]}"; do
  # set up .bashrc
  scp bashrc "$host:/root/.bashrc"
  ssh "$host" mkdir -p $download_dir

  # install crun
  scp ../vendor/crun/* "${host}:${bin_dir}/"
  
  # install conmon
  scp ../vendor/conmon/* "${host}:${bin_dir}/"

  # install crio: scp bundle and run `install`
  scp -r ../vendor/crio/ "${host}:${download_dir}/"
  ssh "$host" << EOF
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

done

