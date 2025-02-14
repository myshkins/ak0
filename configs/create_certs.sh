# create ca key
openssl genrsa -out certs/private/ak0_ca.key 4096

# create ca cert
openssl req -x509 -new -nodes \
  -key certs/private/ak0_ca.key \
  -sha256 -days 36500 \
  -out certs/ak0_ca.crt \
  -subj "/C=US/ST=State/L=City/O=Organization/OU=Unit/CN=AK0-CA"

# prometheus certs
# Create prometheus private key
openssl genrsa -out certs/private/prometheus.key 2048

# Create CSR config
cat << EOF > certs/prometheus-csr.conf
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
prompt = no

[req_distinguished_name]
C = US
ST = State
L = City
O = Organization
OU = Unit
CN = prometheus

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = prometheus
DNS.2 = localhost
IP.1 = 127.0.0.1
EOF

# Generate CSR
openssl req -new -key certs/private/prometheus.key \
  -out certs/prometheus.csr \
  -config certs/prometheus-csr.conf

# Sign the certificate
openssl x509 -req -in certs/prometheus.csr \
  -CA certs/ak0_ca.crt -CAkey certs/private/ak0_ca.key \
  -CAcreateserial -out certs/prometheus.crt \
  -days 36500 -sha256 \
  -extensions v3_req \
  -extfile certs/prometheus-csr.conf

## otel certs
# Create otel private key
openssl genrsa -out certs/private/otel.key 2048

# Create CSR config
cat << EOF > certs/otel-csr.conf
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
prompt = no

[req_distinguished_name]
C = US
ST = State
L = City
O = Organization
OU = Unit
CN = otelcol

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = otelcol
DNS.2 = localhost
IP.1 = 127.0.0.1
EOF

# Generate CSR
openssl req -new -key certs/private/otel.key \
  -out certs/otel.csr \
  -config certs/otel-csr.conf

# Sign the certificate
openssl x509 -req -in certs/otel.csr \
  -CA certs/ak0_ca.crt -CAkey certs/private/ak0_ca.key \
  -CAcreateserial -out certs/otel.crt \
  -days 36500 -sha256 \
  -extensions v3_req \
  -extfile certs/otel-csr.conf
