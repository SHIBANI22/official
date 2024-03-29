# Use the official Golang image as a base image
FROM golang:latest AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY /var/lib/jenkins/workspace/official .

# Build the Go application
RUN go build -o app

# Use a lightweight base image for the final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Expose port 8000
EXPOSE 8000

# Copy the built executable from the build stage
COPY --from=build /app/app .

# Command to run the executable
CMD ["./app"]
