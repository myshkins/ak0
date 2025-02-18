#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path


set -e
echo ""
echo "running create_certs.sh"

org_name="ak0"

rm -rf certs
mkdir -p certs/private

# create ca key and cert
openssl genrsa -out certs/private/ak0_ca.key 4096
openssl req -x509 -new -nodes -key certs/private/ak0_ca.key -sha256 \
  -days 36500 -out certs/ak0_ca.crt -subj "/C=US/O=${org_name}/CN=AK0-CA"

# create config file for prometheus certs
cat > certs/prometheus.cnf <<EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
C = US
O = ${org_name}
CN = prometheus

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = prometheus
DNS.3 = prometheus.local
IP.1 = 127.0.0.1

EOF

# create config file for OpenTelemetry certs
cat > certs/otel.cnf <<EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
C = US
O = ${org_name}
CN = otelcol

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = otelcol
DNS.3 = otelcol.local
IP.1 = 127.0.0.1
EOF


# Create prometheus private key, generate csr, sign cert
openssl genrsa -out certs/private/prometheus.key 4096
echo "generating prom csr"
openssl req -new -key certs/private/prometheus.key -out certs/prometheus.csr \
  -config certs/prometheus.cnf
openssl x509 -req -in certs/prometheus.csr -CA certs/ak0_ca.crt -CAkey certs/private/ak0_ca.key \
  -CAcreateserial -out certs/prometheus.crt -days 36500 -sha256 -extensions req_ext \
  -extfile certs/prometheus.cnf


# Create otel private key, generate csr, sign cert
openssl genrsa -out certs/private/otel.key 4096
openssl req -new -key certs/private/otel.key -out certs/otel.csr \
  -config certs/otel.cnf
openssl x509 -req -in certs/otel.csr -CA certs/ak0_ca.crt -CAkey certs/private/ak0_ca.key \
  -CAcreateserial -out certs/otel.crt -days 36500 -sha256 -extensions req_ext \
  -extfile certs/otel.cnf

# Set proper permissions
chmod 600 certs/private/ak0_ca.key
chmod 644 certs/private/prometheus.key
chmod 644 certs/private/otel.key
chmod 644 certs/ak0_ca.crt
chmod 644 certs/prometheus.crt
chmod 644 certs/otel.crt

echo "certificate generation complete"
echo ""
