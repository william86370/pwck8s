# Go HTTPS Proxy

This project is an HTTPS proxy written in Go. It is designed to handle incoming HTTPS requests, parse client certificate information, and forward the requests to a backend service over HTTP with additional headers containing certificate details.

## Features

- **TLS Support**: Handles incoming HTTPS requests with TLS.
- **Client Certificate Parsing**: Extracts information from client certificates.
- **Header Modification**: Adds client certificate details to the request headers.
- **Reverse Proxy**: Forwards requests to a specified backend service.
- **Environment Variable Configuration**: Configures settings using environment variables.

## Configuration

The service is configured using the following environment variables:

- `PORT`: The port on which the proxy will listen.
- `TLS_CERT`: Path to the TLS certificate.
- `TLS_KEY`: Path to the TLS private key.
- `CA_CERT`: Path to the CA certificate for verifying client certificates.
- `PROXY_URL`: URL of the backend service to proxy to.
- `DEBUG`: Enables debug mode (`true` or `false`).

Additional header configurations:

- `HTTP_HEADER_CN`: The header name for the client's Common Name. Default: `X-Client-Cn`.
- `HTTP_HEADER_DN`: The header name for the client's Distinguished Name. Default: `X-Client-Dn`.

## Prerequisites

- Go 1.x or higher.
- SSL/TLS certificates and keys.

## Setup and Running

1. **Set Environment Variables**: Configure the necessary environment variables as described in the Configuration section.
```bash
PORT=8443
TLS_CERT=/path/to/cert.pem
TLS_KEY=/path/to/cert.key
CA_CERT=/path/to/ca.pem
PROXY_URL=proxy.example.com:8080
DEBUG=false

# HTTP_HEADER_CN=""
# HTTP_HEADER_DN=""
```
2. **Build the Application**:
```bash 
make build-dev
``` 
3. **Run the Application**:
```bash 
make run-dev
```