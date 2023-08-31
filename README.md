# D7024E-Group-3
Repository for the distributed system in the course D7024E at Luleå University of Technology

# Build and run
Build the Docker image running:

```
docker build . -t kadlab
```

Start the network of containers by running:

```
docker-compose up -d   
```

Enter any node by typing:

```
docker exec -it <node-name> sh
```

with an example of a node name being d7024e-group-3-kademliaNodes-1.

This will open a shell environment in the node where we can write different commands. We have installed ping to the image, which makes it possible to message other nodes in the network by typing:

```
ping <node-name>
```