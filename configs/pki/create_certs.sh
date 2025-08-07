#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path


set -e

echo ""
echo "running create_certs.sh"

source ../../.env

rm -f ./certs/*

# create ca.key and ca.crt
openssl genrsa -out certs/ca.key 4096
openssl req -x509 -new -sha512 -nodes -key certs/ca.key -days 3653 -config ca.conf \
  -out certs/ca.crt

certs=(
  "admin"
  "worker-node-0"
  "worker-node-1"
  "kube-proxy"
  "kube-scheduler"
  "kube-controller-manager"
  "kube-api-server"
  "service-accounts"
)

for i in ${certs[*]}; do
  echo "generating ${i} key"
  openssl genrsa -out "certs/${i}.key" 4096

  echo ""
  echo "generating ${i} csr"
  echo "openssl req -new -key certs/${i}.key -sha256 \
    -config ${i}.conf \
    -out certs/${i}.csr"

  openssl req -new -key "certs/${i}.key" -sha256 \
    -config "${i}.conf" \
    -out "certs/${i}.csr"

  echo ""
  echo "generating certs/${i} cert"
  openssl x509 -req -sha256 \
    -days 3653 \
    -CAcreateserial \
    -in "certs/${i}.csr" \
    -CA "certs/ca.crt" \
    -CAkey "certs/ca.key" \
    -extfile "admin.conf" \
    -extensions admin_req_extensions \
    -out "certs/${i}.crt"
done

# chmod 600 certs/private/ak0_ca.key
# chmod 600 certs/ak0_ca.crt

# chown 65534:65534 certs/prometheus.crt

echo ""
echo "yay, certificate generation complete"
