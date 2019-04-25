#!/bin/bash
if [ -z "$WORKSPACE_SHARE_FOLDER" ]; then
  echo '''
On Windows, define WORKSPACE_SHARE_FOLDER=Drive:\\path\\to\\workspace\\ to share folder with the vms (trailing \ is important !)
  '''
  exit 1
fi

export VIRTUALBOX_SHARE_FOLDER=$WORKSPACE_SHARE_FOLDER:workspace
echo "Sharing $VIRTUALBOX_SHARE_FOLDER"

WITH_PROXY=
if [ "$1" = "proxy" ]; then
  WITH_PROXY="--engine-env HTTP_PROXY=$PROXY --engine-env HTTPS_PROXY=$PROXY --engine-env NO_PROXY=$NO_PROXY"
fi

docker-machine create -d "virtualbox" --virtualbox-memory "1024" --virtualbox-cpu-count "1" $WITH_PROXY baby1
docker-machine create -d "virtualbox" --virtualbox-memory "1024" --virtualbox-cpu-count "1" $WITH_PROXY baby2
docker-machine create -d "virtualbox" --virtualbox-memory "1024" --virtualbox-cpu-count "1" $WITH_PROXY baby3