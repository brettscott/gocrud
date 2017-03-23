#!/bin/bash


docker-compose stop web
docker-compose build web

action_local=false
while getopts 'l' flag; do
  case $flag in
    l)
        action_local=true
        ;;
  esac
done

if [ "$action_local" == true ]; then
    docker-compose run --rm --service-ports web ./gocrud
else
    docker-compose run --rm --service-ports web ./build-app.sh
fi