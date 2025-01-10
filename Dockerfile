# Step 1: Use an official Go image as a base image
FROM golang:1.22

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 4: Copy the rest of the application code
COPY . .

# Step 5: Build the Go application
RUN go build -o main .

# Step 6: Expose the port your application listens on (e.g., 8080)
EXPOSE 8080

# Step 7: Command to run the application
CMD ["./main"]
