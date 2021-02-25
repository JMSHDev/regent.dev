#!/bin/bash

# Requires openssl.
#
# This script generates all certificates and keys required for a TLS connection between agent and the platform.
# ca.crt - Certificate Authority (CA) self-signed certificate that that has to be in both platform and the agent.
# server.crt - Platform's certificate (signed by CA) that has to be uploaded to the platform.
# server.key - Platform's private key that has to uploaded to the platform.
#
# Other files generated but not used later on: ca.key - CA private key, ca.srl - list of certificates signed by CA,
# server.crt - platform's certificate,  signing request,
#
# All certificates are saved in platform/emqx/tls where they are automatically uploaded to platform during build time.

 Create self-signed CA certificate.
openssl req -nodes -new -x509  -keyout ca.key -out ca.crt -subj "/C=UK/L=Northampton/O=regenet.dev/CN=regent.dev"

# Next, we need to create a broker's key pair.
openssl genrsa -out server.key 2048

# We use the server key to create a certificate signing request.
openssl req -new -out server.csr -key server.key -subj "/C=UK/L=Northampton/O=regenet.dev/CN=regent.dev"

# Finally we can signed the request with CA certificate.
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 360

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
TLSPATH="$SCRIPTPATH"/../emqx/tls

cp ca.crt "$TLSPATH"/ca.crt
cp server.crt "$TLSPATH"/server.crt
cp server.key "$TLSPATH"/server.key

mkdir -p certificates
mv -t certificates ca.key ca.crt ca.srl server.key server.csr server.crt


