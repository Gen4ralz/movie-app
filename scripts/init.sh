#!/usr/bin/env bash
set -x
set -eo pipefail

# Compile the Protobuf files using 'make compile'
make compile

# Check if the "dev-consul" container is already running
if [[ -n $(docker ps -q -f "name=dev-consul") ]]; then
echo "Container 'dev-consul' is already running."
else
# Run the new container
docker run \
-p 8500:8500 \
-p 8600:8600/udp \
--name=dev-consul \
-d hashicorp/consul agent \
-server \
-ui \
-node=server-1 \
-bootstrap-expect=1 \
-client=0.0.0.0
fi

# Run the Go applications
go run metadata-service/cmd/main.go &
echo "Metadata service started successfully"

go run rating-service/cmd/main.go &
echo "Rating service started successfully"

go run movie-service/cmd/main.go &
echo "Movie service started successfully"

# Wait for all background jobs (services) to finish
wait