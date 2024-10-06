# Use an official Go runtime as a parent image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to leverage Docker cache for dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the entire source code to the working directory
COPY . .

# Build the Go application and place the output in /app/bin
RUN go build -v -o /usr/local/bin/server ./cmd/server

# Copy the docker-entrypoint script and make it executable
COPY ./bin/docker-entrypoint /usr/local/bin/docker-entrypoint
RUN chmod +x /usr/local/bin/docker-entrypoint

# Set the entrypoint for the container
ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]

# Expose the port the server will run on
EXPOSE 8080

# Set the default command to run the server
CMD ["/usr/local/bin/server"]
