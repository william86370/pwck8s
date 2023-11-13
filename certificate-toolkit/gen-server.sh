#!/bin/bash

# Script to generate a new server certificate with X.509 Subject Alternative Names

# Define the server certificate and key file paths
CERT_FILE="server/server.crt"
KEY_FILE="server/server.key"

# Exit if the server certificate and key already exist
if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
    echo "Error: Server certificate and key already exist."
    exit 1
fi

# Create the /server directory if it doesn't exist
mkdir -p server

# Configuration for the server certificate
COUNTRY="US"
ORGANIZATION="L.B. Cloud"
ORG_UNIT="CK8S"
COMMON_NAME="pwck8s.lootbot.cloud"
SAN="DNS:pwck8s.lootbot.cloud,DNS:*.pwck8s.lootbot.cloud,DNS:localhost,IP:127.0.0.1"

# Generate a new RSA 2048 bit server key
openssl genrsa -out "$KEY_FILE" 2048

# Generate a Certificate Signing Request (CSR)
openssl req -new -key "$KEY_FILE" -out server/server.csr \
    -subj "/C=$COUNTRY/O=$ORGANIZATION/OU=$ORG_UNIT/CN=$COMMON_NAME" \
    -addext "subjectAltName = $SAN"

# Check for errors in CSR generation
if [ $? -ne 0 ]; then
    echo "Error: Failed to generate the server CSR."
    exit 1
fi

# Sign the server certificate with the CA certificate and include SAN
openssl x509 -req -in server/server.csr -CA ca/ca.crt -CAkey ca/ca.key -CAcreateserial \
    -out "$CERT_FILE" -days 365 -sha256 \
    -extfile <(printf "subjectAltName = $SAN")

# Check for errors in signing the certificate
if [ $? -ne 0 ]; then
    echo "Error: Failed to sign the server certificate with the CA certificate."
    exit 1
fi

# Display the server certificate information (DN and SAN)
echo "Server Certificate DN:"
openssl x509 -in "$CERT_FILE" -noout -subject

echo "Server Certificate SAN:"
openssl x509 -in "$CERT_FILE" -noout -text | grep -A1 "Subject Alternative Name"