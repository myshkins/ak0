resource "linode_firewall" "k8-node-firewall" {
  label = "k8-node-firewall"

  inbound {
    label    = "allow-ssh-http"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "40020"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }

  inbound {
    label    = "allow-local-http"
    action   = "ACCEPT"
    protocol = "TCP"
    ipv4     = ["10.0.0.0/24"]
  }

  inbound {
    label    = "allow-icmp"
    action   = "ACCEPT"
    protocol = "ICMP"
    ipv4     = ["10.0.0.0/24"]
  }

  inbound_policy = "DROP"

  outbound_policy = "ACCEPT"
}

resource "linode_firewall" "lb-firewall" {
  label = "lb-firewall"

  inbound {
    label    = "allow-ssh-http"
    action   = "ACCEPT"
    protocol = "TCP"
    ports    = "40020, 80, 443"
    ipv4     = ["0.0.0.0/0"]
    ipv6     = ["::/0"]
  }

  inbound {
    label    = "allow-icmp"
    action   = "ACCEPT"
    protocol = "ICMP"
    ipv4     = ["10.0.0.0/24"]
  }

  inbound_policy = "DROP"

  outbound_policy = "ACCEPT"
}

