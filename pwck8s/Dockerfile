# Start with the first stage: Build the application
# Use the official Go image as the base image
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
# This layer is cached, so dependencies will only be re-downloaded if these files change
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Compile the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pwck8s-backend

# Start the second stage: Create the runtime container
FROM scratch

# Add a label to your Docker image (optional)
LABEL maintainer="William Wright <william86370@gmail.com>"
LABEL org.opencontainers.image.source="github.com/william86370/pwck8s.git"

# Copy the compiled binary from the builder stage
COPY --from=builder /app/pwck8s-backend /pwck8s-backend

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./pwck8s-backend"]
