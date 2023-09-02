# D7024E Group 3
Repository for the distributed system in the course D7024E at Luleå University of Technology

# Prerequisites
If you want to run the project locally it's required to download the latest versions of [Docker Desktop](https://www.docker.com/products/docker-desktop/) and [Go](https://go.dev/dl/) (we used go1.21.0). We recommend using VSCode for code editing, but any editor will work. In VSCode we can download extensions for Go, Docker and also Dev Containers which lets us use Docker containers as full-featured development environments.

# Build and run locally
Build the Docker image by running this in a terminal in the D7024E-Group-3 project directory:
```
docker build . -t kadlab
```
Start the network of containers:
```
docker-compose up -d   
```
Enter the shell environment of any node by running:
```
docker exec -it <node-name> sh
```
An example of a node name is d7024e-group-3-kademliaNodes-1.

We have installed the ping command to the image, which makes it possible to message other nodes in the network by typing:
```
ping <node-name>
```

# Deploy to DUST VM
In the future, we would like the procedure of deploying to the DUST VM to be automatic by running a script or similar, but we will do it manually until we have figured out a good way to do it.

## Deploy the source code
Note that all code on the server can be edited using editors such as Vim, but this is bad practice in our case since we don't have any backups of the server. Thus, all code should be edited locally and commited to GitHub, and then deployed to the server using the method described below.

We copy the directory (D7024E-Group-3) containing the source code to the server by using SCP (note that we place the directory in /home/martinaskolin, which is the directory you land in when logging in):
```
scp -r -P 27001 PATH_TO_DIRECTORY\D7024E-Group-3\ martinaskolin@130.240.207.20:/home/martinaskolin/
```
The server is then reached via:
```
ssh -p 27001 martinaskolin@130.240.207.20
```

## Build the image (Dockerfile)
Note that the Docker image only has to be rebuilt if the Dockerfile has changed.

To rebuild an image, go to the D7024E-Group-3 directory on the server and run:
```
sudo docker build . -t kadlab
```
You can check the status of all images by running:
```
sudo docker images
```
To delete the image, run:
```
sudo docker image rm kadlab
```
If there are containers in the stack that need to be updated with the latest image, run:
```
sudo docker service update --image kadlab:latest STACKNAME_kademliaNodes
```

## Deploy Docker stack (docker-compose.yml)
This section refers to starting a Swarm and stack for the first time. If a Swarm and stack is already running, and you have made changes to the docker-compose.yml file that you want to deploy, ignore the rest of this section and simply run the following command in the D7024E-Group-3 directory on the server:
```
sudo docker stack deploy -c docker-compose.yml STACKNAME
```
If you instead want to start a Swarm and stack, follow the steps below:

To start a network of replicas, we need to deploy a stack from the docker-compose.yml file.

The first step to deploying a stack requires that a Swarm is running and that the current node is a Swarm manager, which is initialized with:
```
sudo docker swarm init
```
You can check the status by listing all nodes in a Swarm:
```
sudo docker node ls
```
Then, enter the directory /home/martinaskolin/D7024E-Group-3 which contains the docker-compose.yml file, and deploy the stack by running:
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
To enter containers like we do when running locally, we first need to figure out their names or ID:s. We do that by running this command:
```
sudo docker container ls
```
We can then use either name or ID to enter a container:
```
sudo docker exec -it <NAME/ID> sh
```
From there we can use Unix commands and ping other containers in the network.