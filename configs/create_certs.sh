#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path


set -e

echo ""
echo "running create_certs.sh"

org_name="ak0"
req_config=$(cat <<'END_REQ'
[req]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[req_ext]
subjectAltName = @alt_names

[dn]
C = US
END_REQ
)

rm -rf certs
mkdir -p certs/private

# create ca key and cert
openssl genrsa -out certs/private/ak0_ca.key 4096
openssl req -x509 -new -nodes -key certs/private/ak0_ca.key -sha256 \
  -days 36500 -out certs/ak0_ca.crt -subj "/C=US/O=${org_name}/CN=AK0-CA"

# create config file for prometheus certs
cat > certs/prometheus.cnf <<EOF
$req_config
O = ${org_name}
CN = prometheus

[alt_names]
DNS.1 = localhost
DNS.2 = prometheus
DNS.3 = prometheus.local
DNS.4 = prometheus_ak0
IP.1 = 127.0.0.1
EOF

# create config file for OpenTelemetry certs
cat > certs/otel.cnf <<EOF
$req_config
O = ${org_name}
CN = otelcol

[alt_names]
DNS.1 = localhost
DNS.2 = otelcol
DNS.3 = otelcol.local
DNS.4 = otelcol_ak0
IP.1 = 127.0.0.1
EOF

# create config file for grafana certs
cat > certs/grafana.cnf <<EOF
$req_config
O = ${org_name}
CN = grafana

[alt_names]
DNS.1 = localhost
DNS.2 = grafana
DNS.3 = grafana.local
DNS.4 = grafana_ak0
IP.1 = 127.0.0.1
EOF

# Create prometheus private key, generate csr, sign cert
openssl genrsa -out certs/private/prometheus.key 4096
echo "generating prometheus csr"
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

# Create grafana private key, generate csr, sign cert
openssl genrsa -out certs/private/grafana.key 4096
openssl req -new -key certs/private/grafana.key -out certs/grafana.csr \
  -config certs/grafana.cnf
openssl x509 -req -in certs/grafana.csr -CA certs/ak0_ca.crt -CAkey certs/private/ak0_ca.key \
  -CAcreateserial -out certs/grafana.crt -days 36500 -sha256 -extensions req_ext \
  -extfile certs/grafana.cnf

# todo: set proper permissions, why doesn't this work with the correct perms
chmod 600 certs/private/ak0_ca.key
chmod 600 certs/private/prometheus.key
chmod 600 certs/private/otel.key
chmod 600 certs/private/grafana.key
chmod 666 certs/ak0_ca.crt
chmod 666 certs/prometheus.crt
chmod 666 certs/otel.crt
chmod 666 certs/grafana.crt

chown 10001:10001 certs/private/otel.key
chown 10001:10001 certs/otel.crt
chown 65534:65534 certs/private/prometheus.key
chown 65534:65534 certs/prometheus.crt
chown 472:472 certs/private/grafana.key
chown 472:472 certs/grafana.crt

echo ""
echo "yay, certificate generation complete"
