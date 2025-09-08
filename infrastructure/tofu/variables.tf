variable "linode_api_key" {
  description = "linode api key"
  type        = string
  sensitive   = true
}

variable "nixmo_ssh_pub_key" {
  description = "nixmo ssh pub key"
  type        = string
  sensitive   = false
}

variable "nixabun_ssh_pub_key" {
  description = "nixabun ssh pub key"
  type        = string
  sensitive   = false
}

variable "load_balancer_hostname" {
  description = "load balancer hostname"
  type        = string
  sensitive   = false
}

variable "master_node_0_hostname" {
  description = "master node 0 hostname"
  type        = string
  sensitive   = false
}

variable "worker_node_0_hostname" {
  description = "worker node 0 hostname"
  type        = string
  sensitive   = false
}

variable "worker_node_1_hostname" {
  description = "worker node 1 hostname"
  type        = string
  sensitive   = false
}

