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

# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Install necessary packages
RUN apt-get update && apt-get install -y iputils-ping tmux bash && apt-get clean

# Add the directory where tmux is installed to the PATH
ENV PATH="/usr/bin/tmux:${PATH}"

# Copy the source code from the host into the container
COPY src .

# Set the working directory to be node-cli
WORKDIR /app/cli

# Build an executable for the node-cli
RUN go build -o /app/node-cli

# Set the working directory to be node-cli
WORKDIR /app/kademlia

# Build the Go executable for the Kademlia application
RUN go build -o /app/kademlia-app

# Reset the working directory to /app
WORKDIR /app

# Start a tmux session with the Kademlia application
CMD ["tmux", "new-session", "-s", "kademlia-session", "./kademlia-app"]
