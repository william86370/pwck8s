# Script to generate a new CA certificate and key in the /ca folder 

# check for an existing CA certificate and key
if [ -f ca/ca.crt ] && [ -f ca/ca.key ]; then
    echo "CA certificate and key already exist"
    exit 0
fi

# Create the /ca folder if it doesn't exist
mkdir -p ca
# Generate a new CA certificate and key
# rsa 2048 bit key
# 365 days validity

# Define the default subject for the CA certificate
# /C=Country/O=Organization/OU=Organization-Unit/CN=Common Name
COUNTRY="US"
ORGANIZATION="L.B. Cloud"
ORGANIZATION_UNIT="CK8S"
COMMON_NAME="CK8S PKI Root 1"

# Generate the CA certificate and key
openssl req -x509 -newkey rsa:2048 -days 365 -nodes -keyout ca/ca.key -out ca/ca.crt -subj "/C=$COUNTRY/O=$ORGANIZATION/OU=$ORGANIZATION_UNIT/CN=$COMMON_NAME"

# Check for an Error 
if [ $? -ne 0 ]; then
    echo "Error: Failed to generate the CA certificate and key"
    exit 1
fi

# # Display the CA certificate information DN 
# echo "CA Certificate DN:" 
# openssl x509 -in ca/ca.crt -noout -subject
