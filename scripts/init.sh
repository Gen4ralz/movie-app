#!/usr/bin/env bash
set -x
set -eo pipefail

# Compile the Protobuf files using 'make compile'
make compile

# Check if the dev-consul container exists
if docker ps -a --format '{{.Names}}' | grep -q '^dev-consul$'; then
>&2 echo "The dev-consul container exists, so stop and remove it"
docker stop dev-consul
docker rm dev-consul
fi
# Start the new dev-consul container
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

# Check if the movie_db container exists
if docker ps -a --format '{{.Names}}' | grep -q '^movie_db$'; then
>&2 echo "The movie_db container exists, so stop and remove it"
docker stop movie_db
docker rm movie_db
fi
# Start the new movie_db container
docker run \
-e MYSQL_ROOT_PASSWORD=password \
-e MYSQL_DATABASE=movie \
-p 3306:3306 \
--name=movie_db \
-d mysql:latest

until docker exec movie_db mysql --host=localhost --port=3306 --user=root --password=password --database=movie --execute="SELECT 1"; do
>&2 echo "MySQL is still unavailable - sleeping"
sleep 2
done
>&2 echo "MySQL is up and running on port 3306!"

>&2 echo "Creating schema table..."
docker exec \
-i movie_db mysql movie \
-h localhost \
-P 3306 \
--protocol=tcp \
-uroot \
-ppassword \
< schema/schema.sql

# Run the Go applications
go run metadata-service/cmd/main.go &
>&2 echo "Metadata service started successfully"

go run rating-service/cmd/main.go &
>&2 echo "Rating service started successfully"

go run movie-service/cmd/main.go &
>&2 echo "Movie service started successfully"

# Wait for all background jobs (services) to finish
wait