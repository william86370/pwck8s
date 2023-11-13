# Go HTTP Headers Handler

This Go application is designed to handle HTTP requests by displaying the headers of each incoming request in JSON format. It serves as a simple example of working with HTTP requests in Go.

## Features

- **Request Header Logging**: Logs the headers of each incoming HTTP request.
- **JSON Response**: Converts the request headers to JSON and sends them back in the response.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go 1.x or higher.

### Installing

1. **Clone the repository**:
```bash
git clone https://github.com/william86370/pwck8s.git
```
2. **Navigate to the directory**:
```bash
cd pwck8s/headers
```

### Running the Application

1. **Start the server**:
```bash
make run-dev
```
2. The server will start on `localhost:8080`.

### Usage

- Send any HTTP request to `localhost:8080`, and the server will respond with the headers of the request in JSON format.
- Access `localhost:8080` from a web browser or use a tool like `curl` to make a request.

## Example

Using `curl` to send a request and receive headers:
```bash
curl -i localhost:8080
```

## Built With

- [Go](https://golang.org/) - The Go programming language.