#!/bin/bash

# Prerequisites
# This script expects an arguemnt for the rkm-outpost docker tag

DOCKER_TAG=""
RKM_OUTPOST_GO_BINARY_PATH="./build/package/rkm-outpost/rkm-outpost"

# Check arguments
if [ "$#" -lt 1 ] 
  then
    echo "WARNING: This script expects one argument for the docker tag, you didn't pass one so the script uses the tag: latest"
    DOCKER_TAG="latest"
  elif [ "$#" -gt 1 ]
    then
      echo "ERROR: Your passed too much arguments. We expect only one, the docker tag. Abort..."
      exit 1
  else
      DOCKER_TAG=$1
fi
echo "INFO: Docker Tag is: $DOCKER_TAG"

# Go Build
echo "INFO: Start go build for rkm-outpost"
if [ -f $RKM_OUTPOST_GO_BINARY_PATH ]
  then
    echo "WARNING: rkm-outpost go binary already existed and was deleted by this script to avoid conflicts"
    rm $RKM_OUTPOST_GO_BINARY_PATH
fi
go build -o $RKM_OUTPOST_GO_BINARY_PATH -ldflags="-w -s" ./cmd/rkm-outpost
echo "INFO: Finisched go build for rkm-outpost"

# Docker Build
echo "INFO: Start docker build for rkm-outpost:$DOCKER_TAG"
docker build -t quay.io/onzack/rkm-outpost:$DOCKER_TAG ./build/package/rkm-outpost
echo "INFO: Finisched docker build for rkm-outpost:$DOCKER_TAG"

# Cleanup
echo "INFO: Start cleanup"
rm $RKM_OUTPOST_GO_BINARY_PATH
echo "INFO: Finisched cleanup"