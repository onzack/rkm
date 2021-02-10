#!/bin/bash

# Define varialbes for log output
if [ -f /var/run/secrets/kubernetes.io/serviceaccount/ca.crt ]
  then
    OKLOGTARGET="/proc/1/fd/1"
    ERRORLOGTARGET="/proc/1/fd/2"
  else
    OKLOGTARGET="/dev/stdout"
    ERRORLOGTARGET="/dev/stderr"
fi

# Define global varialbes
KUBECTL_PATH="/usr/local/bin/kubectl"
OPERATOR_SUCCESS="true"
TAG=$(date '+%Y-%m')

echo "Current tag is: $TAG" > $OKLOGTARGET

# Actual script
echo "Patching rkm-outpost cronJob with current tag" > $OKLOGTARGET
$KUBECTL_PATH set -n rkm-outpost image cronjob/rkm-outpost rkm-outpost=$REPOSITORY:$TAG
if (( $? == "0" ))
  then
    echo "OK - patch successful" > $OKLOGTARGET
    exit 0
  else
    echo "ERROR - patch went wrong" > $ERRORLOGTARGET
    exit 1
fi