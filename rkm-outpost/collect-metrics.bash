#!/bin/bash

# Prerequisites
# This script needs this environment variables to be set
# - CLUSTER_NAME
# - INFLUXDB_URL
# - INFLUXDB_PORT
# - INFLUXDB_NAME
# - AUTH_ENABLED
# - INFLUXDB_USER
# - INFLUXDB_PW
# - VERBOSE

STARTTIME=$(date +%s.%N)

# Define varialbes for log output
if [ -f /var/run/secrets/kubernetes.io/serviceaccount/ca.crt ]
  then
    OKLOGTARGET="/proc/1/fd/1"
    ERRORLOGTARGET="/proc/1/fd/2"
  else
    OKLOGTARGET="/dev/stdout"
    ERRORLOGTARGET="/dev/stderr"
fi

# Check environment variables
if [[ -z $CLUSTER_NAME ]]
  then
    echo "ERROR - CLUSTER_NAME environment variable is not set" > $ERRORLOGTARGET
    exit 1
fi

if [[ -z $INFLUXDB_URL ]]
  then
    echo "ERROR - INFLUXDB_URL environment variable is not set" > $ERRORLOGTARGET
    exit 1
fi

if [[ -z $INFLUXDB_PORT ]]
  then
    echo "ERROR - INFLUXDB_PORT environment variable is not set" > $ERRORLOGTARGET
    exit 1
fi

if [[ -z $INFLUXDB_NAME ]]
  then
    echo "ERROR - INFLUXDB_NAME environment variable is not set" > $ERRORLOGTARGET
    exit 1
fi

# Check authentication environment varaibles and define UPLOAD_TO_RKM_MISSION_CONTROL command
if [ $AUTH_ENABLED == "true" ]
  then
    if [[ -z $INFLUXDB_USER ]]
      then
        echo "ERROR - INFLUXDB_USER environment variable is not set" > $ERRORLOGTARGET
        exit 1
    fi
    if [[ -z $INFLUXDB_PW ]]
      then
        echo "ERROR - INFLUXDB_PW environment variable is not set" > $ERRORLOGTARGET
        exit 1
    fi
    UPLOAD_TO_RKM_MISSION_CONTROL () {
      curl -i -XPOST "$INFLUXDB_URL:$INFLUXDB_PORT/write?db=$INFLUXDB_NAME&u=$INFLUXDB_USER&p=$INFLUXDB_PW" --data-binary @$METRICSFILE
    }
  else
    UPLOAD_TO_RKM_MISSION_CONTROL () {
      curl -i -XPOST "$INFLUXDB_URL:$INFLUXDB_PORT/write?db=$INFLUXDB_NAME&u=$INFLUXDB_USER&p=$INFLUXDB_PW" --data-binary @$METRICSFILE
    }
fi

if [[ -z $VERBOSE ]]
  then
    echo "INFO - VERBOSE environment variable is not set, Use VERBOSE=false as default" > $OKLOGTARGET
    VERBOSE="false"
fi

# Define global varialbes
TEMPFILEPATH="/tmp"
METRICSFILE="$TEMPFILEPATH/metrics.txt"
COMPONENTSTATUSESFILE="$TEMPFILEPATH/componentstatuses-metrics.json"
NODESFILE="$TEMPFILEPATH/nodes-metrics.json"
ENDPOINTSFILE="$TEMPFILEPATH/endpoints-metrics.json"
KUBEGUARD_OVERALL_STATUS="1"

# Define global functions
ECHO_DURATION () {
  ENDTIME=$(date +%s.%N)
  DURATION=$(echo "$ENDTIME - $STARTTIME" | bc -l | sed -e 's/^\./0./')
  echo "rkm_script_duration_seconds,cluster=$CLUSTER_NAME value=$DURATION" >> $METRICSFILE
}

ECHO_OVERALL_STATUS () {
  echo "rkm_overall_health,cluster=$CLUSTER_NAME value=$KUBEGUARD_OVERALL_STATUS" >> $METRICSFILE
}

# Remove old temporary files, it they are still present
if [ -f $METRICSFILE ]
  then
    rm $METRICSFILE
fi

if [ -f $COMPONENTSTATUSESFILE ]
  then
    rm $COMPONENTSTATUSESFILE
fi

if [ -f $NODESFILE ]
  then
    rm $NODESFILE
fi

if [ -f $ENDPOINTSFILE ]
  then
    rm $ENDPOINTSFILE
fi

echo "OK - preflight checks successful, start collecting metrics" > $OKLOGTARGET

# Collect metrics
kubectl get endpoints -n default kubernetes >> /dev/null
if (( $? != "0" ))
  then
    echo "rkm_kubeapiserver_health,cluster=$CLUSTER_NAME value=0" >> $METRICSFILE
    KUBEGUARD_OVERALL_STATUS="0"
    ECHO_OVERALL_STATUS
    ECHO_DURATION
    exit 0
  else
    echo "rkm_kubeapiserver_health,cluster=$CLUSTER_NAME value=1" >> $METRICSFILE
fi

kubectl get componentstatuses -o json >> $COMPONENTSTATUSESFILE
kubectl get nodes -o json >> $NODESFILE
kubectl get endpoints -n default kubernetes -o json >> $ENDPOINTSFILE

i="0"
COMPONENTSLENGTH="0"
COMPONENTNAME=""
COMPONENTSTATUS=""
COMPONENTSLENGTH=$(jq '.items | length' $COMPONENTSTATUSESFILE)
while (( $i < $COMPONENTSLENGTH ))
  do
    COMPONENTNAME=$(jq ".items[$i].metadata.name" $COMPONENTSTATUSESFILE)
    COMPONENTSTATUS=$(jq ".items[$i].conditions[0].status" $COMPONENTSTATUSESFILE | tr -d '"')
    if [ $COMPONENTSTATUS == "True" ]
      then
        echo "rkm_component_health,cluster=$CLUSTER_NAME,component=$COMPONENTNAME value=1" >> $METRICSFILE
      else
        KUBEGUARD_OVERALL_STATUS="0"
        echo "rkm_component_health,cluster=$CLUSTER_NAME,component=$COMPONENTNAME value=0" >> $METRICSFILE
    fi
  i=$(($i + 1))
done

i="0"
NODESLENGTH="0"
NODENAME=""
NODESLENGTH=$(jq '.items | length' $NODESFILE)
CONDITIONSLENGTH="0"
COMPONENTSTATUS=""
while (( $i < $NODESLENGTH ))
  do
    i2="0"
    NODENAME=$(jq ".items[$i].metadata.name" $NODESFILE | tr -d '"')
    CONDITIONSLENGTH=$(jq ".items[0].status.conditions | length" $NODESFILE)
    while (( $i2 < $CONDITIONSLENGTH ))
      do
        CONDITIONTYPE=$(jq ".items[$i].status.conditions[$i2].type" $NODESFILE | tr -d '"')
        CONDITIONSTATUS=$(jq ".items[$i].status.conditions[$i2].status" $NODESFILE | tr -d '"')
        case $CONDITIONTYPE in
          Ready)
            if [ $CONDITIONSTATUS != "True" ]
              then
                KUBEGUARD_OVERALL_STATUS="0"
                echo "rkm_node_conditiontype_health,cluster=$CLUSTER_NAME,node=$NODENAME,condition=$CONDITIONTYPE value=0" >> $METRICSFILE
              else
                echo "rkm_node_conditiontype_health,cluster=$CLUSTER_NAME,node=$NODENAME,condition=$CONDITIONTYPE value=1" >> $METRICSFILE
            fi
          ;;
          *)
            if [ $CONDITIONSTATUS != "False" ]
              then
                KUBEGUARD_OVERALL_STATUS="0"
                echo "rkm_node_conditiontype_health,cluster=$CLUSTER_NAME,node=$NODENAME,condition=$CONDITIONTYPE value=0" >> $METRICSFILE
              else
                echo "rkm_node_conditiontype_health,cluster=$CLUSTER_NAME,node=$NODENAME,condition=$CONDITIONTYPE value=1" >> $METRICSFILE
            fi
          ;;
        esac
      i2=$(($i2 + 1))
    done
    i=$(($i + 1))
done

ENDPOINTSCOUNT="0"
ENDPOINTSCOUNT=$(jq '.subsets[0].addresses | length' $ENDPOINTSFILE)
echo "rkm_apiserver_endpoints_total,cluster=$CLUSTER_NAME value=$ENDPOINTSCOUNT" >> $METRICSFILE

ECHO_OVERALL_STATUS
ECHO_DURATION

# Upload metrics to RKM Mission Control
echo "OK - collecting metrics successful, start uploading to RKM mission control: $INFLUXDB_URL:$INFLUXDB_PORT/write?db=$INFLUXDB_NAME" > $OKLOGTARGET

if [ $VERBOSE == "true" ]
  then
    echo "Print Metricsfile as VERBOSE=true"
    cat $METRICSFILE > $OKLOGTARGET
fi

UPLOAD_TO_RKM_MISSION_CONTROL

if (( $? != "0" ))
  then
    echo "ERROR - uploading to RKM mission control not successful" > $ERRORLOGTARGET
    exit 1
  else
    echo "OK - uploading to RKM mission control successful" > $OKLOGTARGET
fi

exit 0