# Script to generate a new client certificate and key in the /client folder

# check for an existing client certificate and key
if [ -f client/client.crt ] && [ -f client/client.key ]; then
    echo "Error: Client certificate and key already exist"
    exit 1
fi

# Create the /client folder if it doesn't exist
mkdir -p client

# Generate a new client certificate and key
# rsa 2048 bit key
# 365 days validity

# Define the default subject for the client certificate
# /C=Country/O=Organization/OU=Organization-Unit/OU=Organization-Unit/OU=Organization-Unit/OU=Organization-Unit/CN=Common Name
COUNTRY="US"
ORGANIZATION="L.B. Cloud"
ORGANIZATION_UNIT="CK8S"
ORGANIZATION_UNIT2="LB"
ORGANIZATION_UNIT3="D002"
COMMON_NAME="Dimmadome Doug Test dtdimma"
EMAIL="dtdimma@lootbot.cloud"

# Generate the client certificate and key  with email as a subject alternative name
openssl req -newkey rsa:2048 -days 365 -nodes -keyout client/client.key -out client/client.csr -subj "/C=$COUNTRY/O=$ORGANIZATION/OU=$ORGANIZATION_UNIT/OU=$ORGANIZATION_UNIT2/OU=$ORGANIZATION_UNIT3/CN=$COMMON_NAME" -addext "subjectAltName = email:$EMAIL"
# Check for an Error
if [ $? -ne 0 ]; then
    echo "Error: Failed to generate the client certificate and key"
    exit 1
fi

# Sign the client certificate with the CA certificate
openssl x509 -req -in client/client.csr -CA ca/ca.crt -CAkey ca/ca.key -CAcreateserial -out client/client.crt -days 365

# Check for an Error
if [ $? -ne 0 ]; then
    echo "Error: Failed to sign the client certificate with the CA certificate"
    exit 1
fi

# Display the client certificate information DN
echo "Client Certificate DN:"
openssl x509 -in client/client.crt -noout -subject
