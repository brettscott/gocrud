#!/bin/bash


build_containers_locally=true
stop_containers=true
is_express_mode=false
action_local=false
action_ssh=false
while getopts 'lx' flag; do
  case $flag in
    l)
        action_local=true
        ;;
    x)
        build_containers_locally=false
        stop_containers=false
        is_express_mode=true
        ;;
    y)
        action_ssh=true
        ;;
  esac
done

if [ "$build_containers_locally" == true ]; then
    docker-compose stop
    docker-compose rm -f
    docker-compose build
fi

if [ "$action_local" == true ]; then
    docker-compose run --rm --service-ports web ./build-locally.sh
elif [ "$action_ssh" == true ]; then
    docker-compose run --rm --service-ports web
else
    docker-compose run --rm --service-ports web ./build-app.sh
fi

if [ "$stop_containers" == true ]; then
    docker-compose stop
    docker-compose rm -f
fi
