# D7024E Group 3
Repository for the distributed system in the course D7024E at Luleå University of Technology

# Prerequisites
If you want to run the project locally it's required to download the latest versions of [Docker Desktop](https://www.docker.com/products/docker-desktop/) and [Go](https://go.dev/dl/) (we used go1.21.0). We recommend using VSCode for code editing, but any editor will work. In VSCode we can download extensions for Go, Docker and also Dev Containers which lets us use Docker containers as full-featured development environments.

# Build and run locally
Build the Docker image by running this in a terminal in the D7024E-Group-3 project directory:
```
docker build . -t kadlab
```
Start the network of containers by running:
```
docker-compose up -d   
```
Enter the shell environment of any node by typing:
```
docker exec -it <node-name> sh
```
An example of a node name is d7024e-group-3-kademliaNodes-1.

We have installed the ping command to the image, which makes it possible to message other nodes in the network by typing:
```
ping <node-name>
```

# Deploy to DUST VM
In the future, we would like this process to be automatic by running a script or similar, but we will do it manually until we have figured out the best deployment procedure.

## Deploy the image
First, export the Docker image as a tarball file:
```
docker save -o kadlab_image.tar kadlab
```
Copy the tarball file to your remote server using SCP (Secure Copy) and then enter the password:
```
scp -P 27001 kadlab_image.tar martinaskolin@130.240.207.20:/home/martinaskolin/
```
SSH to the server, which will place you in the /home/martinaskolin directory:
```
ssh -p 27001 martinaskolin@130.240.207.20
```
Finally, import the Docker image on the remote server:
```
sudo docker load -i kadlab_image.tar
```
You can check the status of all images by running:
```
sudo docker images
```
To delete the image, run:
```
sudo docker image rm kadlab
```

## Deploy the source code
All source code must be present on the VM in order for us to run it. Add the directory containing the source code to the server by using SCP:
```
scp -r -P 27001 .\D7024E-Group-3\ martinaskolin@130.240.207.20:/home/martinaskolin/
```

## Deploy Docker stack
To start our network of replicas, we need to deploy a stack from the docker-compose.yml file.

This first requires that a Swarm is running and that the current node is a Swarm manager, which is initialized with:
```
sudo docker swarm init
```
To then list all nodes in a swarm, run:
```
sudo docker node ls
```
Enter the directory /home/martinaskolin/D7024E-Group-3 which contains the docker-compose.yml file, and deploy the stack by running:
```
sudo docker stack deploy -c docker-compose.yml STACKNAME
```
This will create a network called STACKNAME_kademlia_network and a service called STACKNAME_kademliaNodes (replace STACKNAME with an appropriate name). The replicas-variable in docker-compose.yml states how many containers will be created.

To list all containers in a stack, run:
```
sudo docker stack ps STACKNAME
```

You can list and remove stacks, services and networks with the following commands:
```
sudo docker stack ls
sudo docker stack rm STACKNAME

sudo docker service ls
sudo docker service rm SERVICENAME

sudo docker network ls
sudo docker network rm NETWORKNAME
```

## Enter a container
To enter containers like we do when running locally, we first need to figure out their names or ID:s. We do that by entering:
```
sudo docker container ls
```
We can then use either name or ID to enter a container:
```
sudo docker exec -it <NAME/ID> sh
```
From there we can use Unix commands and ping other containers in the network.