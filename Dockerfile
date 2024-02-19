# Stage 1: Build the application
FROM golang:latest AS build

WORKDIR /go/src/app

COPY . .

# Build the Go application
RUN go get -d -v ./...
RUN go install -v ./...

# Stage 2: Create a minimal image
FROM alpine:latest

WORKDIR /app

# Expose port 8000
EXPOSE 8000

# Copy the built executable from the first stage
COPY --from=build /go/bin/app .

# Command to run the executable
CMD ["./app"]
