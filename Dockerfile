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
RUN apt-get update

# Copy the source code from the host into the container
COPY src .

# Build an executable for the kademlia-app
RUN go build -o kademlia-app

# Start the Kademlia application
CMD ["./kademlia-app"]
