# D7024E Group 3
Repository for the distributed system in the course D7024E at Luleå University of Technology.

# Prerequisites
If you want to run the project locally it's required to download the latest versions of [Docker Desktop](https://www.docker.com/products/docker-desktop/) and [Go](https://go.dev/dl/) (we used go1.21.0). We recommend using VSCode for code editing, but any editor will work. In VSCode we can download extensions for Go, Docker and also Dev Containers which lets us use Docker containers as full-featured development environments.

# Build and run locally
Build the Docker image by running this in a terminal in the D7024E-Group-3 project directory:
```
docker build . -t kadlab
```
Start the network of containers:
```
docker compose up -d   
```
Enter the shell environment of any node by running:
```
docker exec -it <node-name> sh
```
To enter the Go program directly, we instead run:
```
docker attach <node-name>
```
To see the logs of a node:
```
docker logs <node-name>
```

An example of a node name is d7024e-group-3-kademliaNodes-1.

# Deploy to DUST VM
Any pushes to `main`, either directly or via pull requests, will result in an automatic deployment to the DUST VM. The deployment is performed by a GitHub Action (see `.github/workflows/main.yml`), which builds the Docker image and deploys the Docker containers accoring to the `docker-compose.yml` file.

The code can also be copied to the server manually via SSH:
```
scp -r -P 27001 PATH_TO_DIRECTORY\D7024E-Group-3\ martinaskolin@130.240.207.20:/home/martinaskolin/
```
The server is then reached via:
```
ssh -p 27001 martinaskolin@130.240.207.20
```
All Docker commands that are available locally, can also be executed on the server. Before starting new containers, the old ones should be removed. Same goes for the image.

Commands flags such as `ls` and `rm` are useful for listing and removing entities, such as listing all existing containers:
```
sudo docker container ls
```

# Generate HTML Coverage Report

## Run test with coverage
First, enter the `src` directory:
```
cd PATH_TO_DIRECTORY\D7024E-Group-3\src
```
To run your tests with coverage analysis, use the `go test` command with the `-cover flag`:
```
go test -cover ./... -coverprofile=c.out
```
The `./...` argument tells Go to recursively run tests in all packages within the current directory and its subdirectories.

Then run:
```
go tool cover -html=c -o coverage.html
```
or
```
go tool cover -html=c
```
to generate an HTML coverage report to visualize the coverage in more detail.
