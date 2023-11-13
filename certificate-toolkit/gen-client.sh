#!/bin/bash

# Script to generate a new client certificate and key with an email as SAN

# Define the client certificate and key file paths
CERT_FILE="client/client.crt"
KEY_FILE="client/client.key"

# Exit if the client certificate and key already exist
if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
    echo "Client certificate and key already exist."
    exit 0
fi

# Create the /client directory if it doesn't exist
mkdir -p client

# Configuration for the client certificate
COUNTRY="US"
ORGANIZATION="L.B. Cloud"
ORG_UNIT="CK8S"
COMMON_NAME="Dimmadome Doug Test dtdimma"
EMAIL="dtdimma@lootbot.cloud"
SAN="email:$EMAIL"

# Generate a new RSA 2048 bit client key
openssl genrsa -out "$KEY_FILE" 2048

# Generate a Certificate Signing Request (CSR) with email as SAN
openssl req -new -key "$KEY_FILE" -out client/client.csr \
    -subj "/C=$COUNTRY/O=$ORGANIZATION/OU=$ORG_UNIT/CN=$COMMON_NAME" \
    -addext "subjectAltName = $SAN"

# Check for errors in CSR generation
if [ $? -ne 0 ]; then
    echo "Error: Failed to generate the client CSR."
    exit 1
fi

# Sign the client certificate with the CA certificate
openssl x509 -req -in client/client.csr -CA ca/ca.crt -CAkey ca/ca.key -CAcreateserial \
    -out "$CERT_FILE" -days 365 -sha256 \
    -extfile <(printf "subjectAltName = $SAN")

# Check for errors in signing the certificate
if [ $? -ne 0 ]; then
    echo "Error: Failed to sign the client certificate with the CA certificate."
    exit 1
fi

# # Display the client certificate information (DN)
# echo "Client Certificate DN:"
# openssl x509 -in "$CERT_FILE" -noout -subject

# # Optionally, display the client certificate SAN (email)
# echo "Client Certificate SAN:"
# openssl x509 -in "$CERT_FILE" -noout -text | grep -A1 "Subject Alternative Name"
