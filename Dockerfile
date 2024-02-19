# Stage 1: Build the application
FROM golang:latest AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go application
RUN go build -o app

# Stage 2: Create a minimal image
FROM alpine:latest

WORKDIR /app

# Expose port 8000
EXPOSE 8000

# Copy the built executable from the previous stage
COPY --from=build /app/app .

# Command to run the executable
CMD ["./app"]
