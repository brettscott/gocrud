#!/bin/bash -e


build_containers_locally=true
stop_containers=true
is_express_mode=false
action_local=false
action_ssh=false
verbose=false
while getopts 'lxy' flag; do
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

function on_error {
  if [ "$verbose" == true ]; then
      echo "Displaying docker compose logs:"
      docker-compose logs
  else
    echo "  NOTICE: Docker error logs can be displayed here by adding -v"
  fi
  stopContainers
  exit 1
}
trap on_error ERR

function ctrl_c() {
    printf "\nCTRL-C detected in run.sh\n"
    stopContainers
}
trap ctrl_c ERR INT SIGHUP SIGINT SIGTERM


function stopContainers() {
    docker-compose stop
    docker-compose rm -f
    echo "  "
    echo " >  Stopped containers ..."
    echo "  "
}

if [ "$build_containers_locally" == true ]; then
    stopContainers
    docker-compose build

    echo "  "
    echo " >  Containers built ..."
    echo "  "
fi

if [ "$action_local" == true ]; then
    docker-compose run --rm --service-ports web
elif [ "$action_ssh" == true ]; then
    docker-compose run --rm --service-ports web bash
else
    docker-compose run --rm --service-ports web ./build-app.sh
fi

if [ "$stop_containers" == true ]; then
    stopContainers
fi

echo " "
docker-compose ps
echo " "
echo " >  End ..."
echo " "