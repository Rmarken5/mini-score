# Use a more recent version of the Go base image
FROM golang:1.20-alpine AS build

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files separately to leverage Docker layer caching
COPY ../go.mod .
COPY ../go.sum .

# Download Go module dependencies
RUN go mod download

# Copy the rest of the source code
COPY .. .

# Build the application
RUN go build -o mini-score ./service/cmd/scheduler/main.go

# Use a smaller base image for the final runtime image
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=build /app/mini-score .

# Set the command to run the application
CMD ["./mini-score"]
