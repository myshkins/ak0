terraform {
  required_providers {
    linode = {
      source  = "linode/linode"
      version = "3.0.0"
    }
  }
}

provider "linode" {
  token = var.linode_api_key
}

resource "linode_vpc" "ak0-vpc" {
  label = "ak0-vpc"
  region = "us-lax"
}

resource "linode_vpc_subnet" "ak0-vpc-subnet" {
  vpc_id = linode_vpc.ak0-vpc.id
  label = "ak0-vpc-subnet"
  ipv4 = "10.0.0.0/24"
}

resource "random_password" "temp-password" {
  for_each = local.instances
  length  = 32
  special = true
}

resource "random_password" "lb-password" {
  length  = 32
  special = true
}

locals {
  instances = {
    (var.load_balancer_hostname) = {
      label = "ak0-load-balancer"
      machine_type = "g6-nanode-1"
      hostname = var.load_balancer_hostname
      fqdn = ""
      public_interface_primary = true
      firewall_id = linode_firewall.lb-firewall.id
      vpc_ip = "10.0.0.2"
    }
    (var.master_node_0_hostname) = {
      label = "ak0-master-node-0"
      machine_type = "g6-standard-1"
      hostname = var.master_node_0_hostname
      fqdn = "master_node_0.kubernetes.local"
      public_interface_primary = true
      firewall_id = linode_firewall.k8-node-firewall.id
      vpc_ip = "10.0.0.3"
    }
    (var.worker_node_0_hostname) = {
      label = "ak0-worker-node-0"
      machine_type = "g6-nanode-1"
      hostname = var.worker_node_0_hostname
      fqdn = "worker_node_0.kubernetes.local"
      public_interface_primary = true
      firewall_id = linode_firewall.k8-node-firewall.id
      vpc_ip = "10.0.0.4"
    }
    (var.worker_node_1_hostname) = {
      label = "ak0-worker-node-1"
      machine_type = "g6-nanode-1"
      hostname = var.worker_node_1_hostname
      fqdn = "worker_node_1.kubernetes.local"
      public_interface_primary = true
      firewall_id = linode_firewall.k8-node-firewall.id
      vpc_ip = "10.0.0.5"
    }
  }
}

locals {
  # Generate hosts file entries
  hosts_entries = [
    for key, value in local.instances:
    "${value.vpc_ip} ${value.fqdn} ${value.hostname}"
    if key != (var.load_balancer_hostname)
  ]

  # Complete hosts file content
  host_file_content = join("\n", local.hosts_entries)
}

output "host_file_content" {
  value = local.host_file_content
}

output "linode_ips" {
  value = {
    for key, vm in linode_instance.ak0-vm : key => vm.ip_address
  }
  sensitive = true
}

resource "linode_stackscript" "ak0-vm-init-script" {
  label = "ak0-vm-init-script"
  description = "Initializes a vm with proper hostname"
  script = <<EOF
#!/bin/bash
# <UDF name="hostname" label="hostname" default="ak0">
# <UDF name="host_file_content" label="host_file_content" default="added by tofu">
touch /root/ak0up
apt-get -q update && apt-get -q -y upgrade
apt-get install kitty-terminfo
hostnamectl set-hostname "$HOSTNAME"
echo -e "\nadded by tofu\n$HOST_FILE_CONTENT" >> /etc/hosts
# sed -i -e 's/^#Port 22$/Port 40020/' -e 's/^PasswordAuthentication yes$/PasswordAuthentication no/' /etc/ssh/sshd_config
# systemctl restart ssh
EOF
  images = ["linode/debian12"]
  rev_note = "initial version"
}

resource "linode_instance" "ak0-vm" {
  for_each = local.instances
  label  = each.value.label
  image  = "linode/debian12"
  region = "us-lax"
  type   = each.value.machine_type
  authorized_keys = [var.nixmo_ssh_pub_key, var.nixabun_ssh_pub_key]
  root_pass  = random_password.temp-password[each.key].result
  tags       = ["k8", "worker-node"]
  private_ip = true
  firewall_id = each.value.firewall_id
  stackscript_id = linode_stackscript.ak0-vm-init-script.id
  stackscript_data = {
    "hostname" = each.value.hostname
    "host_file_content" = local.host_file_content
  }

  interface {
    purpose   = "public"
    subnet_id = linode_vpc_subnet.ak0-vpc-subnet.id
    primary = each.value.public_interface_primary
  }

  interface {
    purpose   = "vpc"
    subnet_id = linode_vpc_subnet.ak0-vpc-subnet.id
    primary = !each.value.public_interface_primary
    ipv4 {
      vpc = each.value.vpc_ip
    }
  }
}

