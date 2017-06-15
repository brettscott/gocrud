#!/bin/bash -e


build_containers_locally=true
stop_containers=true
is_express_mode=false
action_local=false
action_ssh=false
verbose=false
while getopts 'lvxy' flag; do
  case $flag in
    l)
        action_local=true
        ;;
    v)
        verbose=true
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

function ctrl_c() {
    if [ "$verbose" == true ]; then
        echo " "
        echo " "
        echo " >  Displaying docker compose logs:"
        docker-compose logs
    else
        echo " "
        echo " "
        echo "  >  NOTICE: Docker error logs can be displayed here by adding -v"
    fi
    stopContainers
}
trap ctrl_c ERR INT SIGHUP SIGINT SIGTERM


function stopContainers() {
    if [ "$stop_containers" == true ]; then
        docker-compose stop
        docker-compose rm -f
        printf "\n >  Stopped containers ...\n\n"
    fi
}

if [ "$build_containers_locally" == true ]; then
    stopContainers
    docker-compose build
    printf "\n >  Containers built ...\n\n"
else
    printf "\n >  Not building container ...\n\n"
fi

if [ "$action_local" == true ]; then
    if [ "$build_containers_locally" == false ]; then
        printf "\n >  Running \"go install\" to ensure latest changes are in binary ...\n\n"
        docker-compose run --rm --service-ports web ./build-run.sh
    else
        docker-compose run --rm --service-ports web
    fi
elif [ "$action_ssh" == true ]; then
    docker-compose run --rm --service-ports web bash
else
    docker-compose run --rm --service-ports web ./build-tests.sh
fi

if [ "$stop_containers" == true ]; then
    stopContainers
fi

echo " "
docker-compose ps
printf "\n >  End ...\n\n"
