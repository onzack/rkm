# RKM - Remote Kubernetes Monitoring
Simple monitoring tool for local and remote Kubernetes clusters.  
Used technonlogies:
* Docker
* Bash
* InfluxDB
* Grafana

# Architecture
![A Sample Graph for visualization ](./Docs/rkm.png)

# Grafana Dashbaord
![Grafana Dashboard ](https://github.com/onzack/rkm/blob/main/Docs/rkm-mission-control.png)

# Installation

## RKM Mission Control

### Install InfluxDB
The RKM-Outpost metrics are stored in an InfluxDB.  
If you don't already have an InfluxDB, install one for example with the [InfluxDB Helm Chart](https://github.com/influxdata/helm-charts/tree/master/charts/influxdb) on Kubernetes.  
The InfluxDB and RKM-Outpost have a few dependencies:
1. The InfluxDB must be reachable for RKM-Outpost
3. An InfluxDB database for the RKM-Outpost metrics, for example "rkm"
4. User and password for the InfluxDB database

### Install Grafana
If your don't have a Grafana instance, install one for example with the [Grafana Helm Chart](https://github.com/helm/charts/tree/master/stable/grafana) on Kubernetes. Grafana must have access to the InfluxDB database.  
Add the InfluxDB database as a datasource to Grafana and import the [RKM Mission Control Dashboard](https://github.com/onzack/rkm/blob/main/Grafana/rkm-mission-control-dashboard.json).

## RKM Outpost
Install RKM Outpost Helm Chart:  
1. Create rkm-outpost namespace:  
`kubectl create namespace rkm-outpost`
2. Create secret for InfluxDB authentication:  
`kubectl create secret generic rkm-secrets --from-literal=INFLUXDB_USER=<user> --from-literal=INFLUXDB_PW=<password> -n rkm-outpost`
3. Clone this repository:  
`git clone https://github.com/onzack/rkm.git`
4. Create and adjust the custom-values.yaml file for your cluster:  
`cp ./rkm/custom-values.yaml cluster-1-values.yaml`  
`vim custom-values.yaml`  

**Importat**
Take a special look at the proxy and CA configs, if you run RKM-Outpost behind a enterprise proxy or the InfluxDB SSL/TLS certificate is signed by a not well known CA.
And keep in mind, that the proxy support is early beta and not well tested at the moment.

5. Install Helm Chart:  
`helm install -f custom-values.yaml -n rkm-outpost rkm-outpost ./rkm/Helm/rkm-outpost`  

# Docker repositories
RKM-Outpost: https://quay.io/repository/onzack/rkm-outpost  
RKM-Outpost operator: https://quay.io/repository/onzack/rkm-outpost-operator

# Licence
Copyright 2021 ONZACK AG

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
