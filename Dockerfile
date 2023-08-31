# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab

# Install a small os
FROM alpine:latest

# Install necessary packages
RUN apk add --no-cache iputils
RUN apk add --no-cache bash

# Use the official Go image as the base image
#FROM golang:latest

# Set the working directory inside the container
#WORKDIR /app

# Copy the go.mod file to the container
#COPY go.mod .

# Copy the main file from the host into the container. We only need the main file now when we are testing
#COPY src/main.go .

# Build the Go executable
#RUN go build -o app

# Specify the command to run when the container starts
#CMD ["./app"]
CMD bash