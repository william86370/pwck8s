# Start with the first stage: Build the application
# Use the official Go image as the base image
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
# This layer is cached, so dependencies will only be re-downloaded if these files change
COPY go.mod ./

# Download the dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Compile the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o x509-proxy

# Start the second stage: Create the runtime container
FROM scratch

# Add a label to your Docker image (optional)
LABEL maintainer="William Wright <william86370@gmail.com>"
LABEL description="A simple x509 proxy server"

# Copy the compiled binary from the builder stage
COPY --from=builder /app/x509-proxy /x509-proxy

# Expose port 8443
EXPOSE 8443
ENV PORT=8443

# Command to run the executable
CMD ["./x509-proxy"]
