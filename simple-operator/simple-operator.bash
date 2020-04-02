#!/bin/bash

export IFS=";"

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

for COMPONENT in $COMPONENTS; do
  # Actual script
  echo "Patching ${COMPONENT} cronJob with current tag" > $OKLOGTARGET
  $KUBECTL_PATH set -n ${K8S_NAMESPACE} image ${COMPONENT} ${COMPONENT/*\//}=$REPOSITORY:$TAG
  if (( $? == "0" ))
    then
      echo "OK - patch successful" > $OKLOGTARGET
      exit 0
    else
      echo "ERROR - patch went wrong" > $ERRORLOGTARGET
      exit 1
  fi
done