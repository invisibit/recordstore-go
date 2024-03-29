# Use an official Golang runtime as a parent image
FROM golang:1.21

# Set the working directory to /go/src/app
WORKDIR /go/src/recordstore-go

# Copy the local package files to the container's workspace
COPY . .
COPY default.env .env

# Build the Go application
RUN go get -d -v ./...
RUN go install -v ./...
#RUN go build

# Set the entry point to your application
CMD ["recordstore-go"]

# Document that the service listens on port 8080
EXPOSE 8080
