terraform {
  required_providers {
    linode = {
      source  = "linode/linode"
      version = "3.0.0"
    }
  }
}

provider "linode" {
  token = var.api_key
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
    "load-balancer" = {
      label = "ak0-load-balancer",
      machine_type = "g6-nanode-1"
      hostname = "ak0_lb"
      public_interface_primary = true
      firewall_id = linode_firewall.lb-firewall.id
      vpc_ip = "10.0.0.2"
    }
    "master" = {
      label = "ak0-master-node-0",
      machine_type = "g6-standard-1"
      hostname = "prismo"
      public_interface_primary = true
      firewall_id = linode_firewall.k8-node-firewall.id
      vpc_ip = "10.0.0.3"
    }
    "worker0" = {
      label = "ak0-worker-node-0"
      machine_type = "g6-nanode-1"
      hostname = "gunter0"
      public_interface_primary = true
      firewall_id = linode_firewall.k8-node-firewall.id
      vpc_ip = "10.0.0.4"
    }
    "worker1" = {
      label = "ak0-worker-node-1"
      machine_type = "g6-nanode-1"
      hostname = "gunter1"
      public_interface_primary = true
      firewall_id = linode_firewall.k8-node-firewall.id
      vpc_ip = "10.0.0.5"
    }
  }
}

resource "linode_stackscript" "ak0-vm-init-script" {
  label = "ak0-vm-init-script"
  description = "Initializes a vm with proper hostname"
  script = <<EOF
#!/bin/bash
# <UDF name="hostname" label="hostname" default="ak0">
apt-get -q update && apt-get -q -y upgrade
apt-get install kitty-terminfo
hostnamectl set-hostname $HOSTNAME
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
  authorized_keys = var.ssh_pub_keys
  root_pass  = random_password.temp-password[each.key].result
  tags       = ["k8", "worker-node"]
  private_ip = true
  firewall_id = each.value.firewall_id
  stackscript_id = linode_stackscript.ak0-vm-init-script.id
  stackscript_data = {
    "hostname" = each.value.hostname
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

