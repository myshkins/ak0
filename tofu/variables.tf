variable "api_key" {
  description = "linode api key"
  type        = string
  sensitive   = true
}

variable "ssh_pub_keys" {
  description = "ssh pub keys"
  type = list(string)
  sensitive = false
}


