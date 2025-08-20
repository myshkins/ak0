#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")

set -e

cd "${script_path}/../tofu"

source ../../.env

echo ""
echo "running configure_ssh.sh"

ssh_config_path=/home/myshkins/.ssh/config

nodes=(
    "$MASTER_NODE_0_HOSTNAME"
    "$WORKER_NODE_0_HOSTNAME" 
    "$WORKER_NODE_1_HOSTNAME"
    "$LOAD_BALANCER_HOSTNAME"
)

entries=""

for host in "${nodes[@]}"; do
  ip_addr=$(tofu output -json linode_ips | jq -r ".$host")
  entries+=$(cat <<EOF
Host $host
  Hostname $ip_addr
  Port 40020
  User root
  IdentityFile /home/myshkins/.ssh/key
EOF
)$'\n'$'\n'
done

sed -i.bak '/## START AK0 CONFIGS ##/,/## END AK0 CONFIGS ##/{/## START AK0 CONFIGS ##/!{/## END AK0 CONFIGS ##/!d}}' $ssh_config_path
printf "%s\n" "$entries" | sed -ie "/## START AK0 CONFIGS ##/r /dev/stdin" $ssh_config_path

# todo: add vms to known hosts
