#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")

set -e

cd "${script_path}/../configs/pki"

source ../../.env

echo ""
echo "running create_certs.sh"


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

for i in "${certs[@]}"; do
  if [[ -f "${i}.template" ]]; then
    echo "writing conf for ${i}"
    envsubst < "${i}.template" > "${i}.conf"
  fi
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
    -extfile "${i}.conf" \
    -extensions "${i}_req_extensions" \
    -out "certs/${i}.crt"
done

chmod 600 certs/*.key
chmod 644 certs/*.crt


echo ""
echo "yay, certificate generation complete"
