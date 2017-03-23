#!/bin/bash


docker-compose stop web
docker-compose run --rm --service-ports web ./build-app.sh