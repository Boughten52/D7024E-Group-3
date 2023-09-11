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
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the source code from the host into the container
COPY src .

# Build the Go executable (we call it node-cli)
RUN go build -o node-cli

# We start a bash shell when the container starts
CMD bash