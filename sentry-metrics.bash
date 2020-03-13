#!/bin/bash

STARTTIME=$(date +%s.%N)

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
  echo "sentry_script_duration_seconds,cluster=$CLUSTER_NAME value=$DURATION" >> $METRICSFILE
}

ECHO_OVERALL_STATUS () {
  echo "sentry_overall_health,cluster=$CLUSTER_NAME value=$KUBEGUARD_OVERALL_STATUS" >> $METRICSFILE
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

# Actual script
kubectl get endpoints -n default kubernetes >> /dev/null
if (( $? != "0" ))
  then
    echo "sentry_kubeapiserver_health,cluster=$CLUSTER_NAME value=0" >> $METRICSFILE
    KUBEGUARD_OVERALL_STATUS="0"
    ECHO_OVERALL_STATUS
    ECHO_DURATION
    exit 0
  else
    echo "sentry_kubeapiserver_health,cluster=$CLUSTER_NAME value=1" >> $METRICSFILE
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
        echo "sentry_component_health,cluster=$CLUSTER_NAME,component=$COMPONENTNAME value=1" >> $METRICSFILE
      else
        KUBEGUARD_OVERALL_STATUS="0"
        echo "sentry_component_health,cluster=$CLUSTER_NAME,component=$COMPONENTNAME value=0" >> $METRICSFILE
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
                echo "sentry_node_conditiontype_health,cluster=$CLUSTER_NAME,node=$NODENAME,condition=$CONDITIONTYPE value=0" >> $METRICSFILE
              else
                echo "sentry_node_conditiontype_health,cluster=$CLUSTER_NAME,node=$NODENAME,condition=$CONDITIONTYPE value=1" >> $METRICSFILE
            fi
          ;;
          *)
            if [ $CONDITIONSTATUS != "False" ]
              then
                KUBEGUARD_OVERALL_STATUS="0"
                echo "sentry_node_conditiontype_health,cluster=$CLUSTER_NAME,node=$NODENAME,condition=$CONDITIONTYPE value=0" >> $METRICSFILE
              else
                echo "sentry_node_conditiontype_health,cluster=$CLUSTER_NAME,node=$NODENAME,condition=$CONDITIONTYPE value=1" >> $METRICSFILE
            fi
          ;;
        esac
      i2=$(($i2 + 1))
    done
    i=$(($i + 1))
done

ENDPOINTSCOUNT="0"
ENDPOINTSCOUNT=$(jq '.subsets[0].addresses | length' $ENDPOINTSFILE)
echo "sentry_apiserver_endpoints_total,cluster=$CLUSTER_NAME value=$ENDPOINTSCOUNT" >> $METRICSFILE

ECHO_OVERALL_STATUS
ECHO_DURATION

curl -i -XPOST "$INFLUXDB_URL/write?db=$INFLUXDB_NAME&u=$INFLUXDB_USER&p=$INFLUXDB_PW" --data-binary @$METRICSFILE

exit 0