# Stage 1: Build the application
FROM golang:latest AS build

WORKDIR /app

COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Stage 2: Create a minimal image
FROM alpine:latest

WORKDIR /app

# Expose port 8000
EXPOSE 8000

# Copy the built executable from the first stage
COPY --from=build /app/app .

# Command to run the executable
CMD ["./app"]
