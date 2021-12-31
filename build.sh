#!/bin/bash
 
# available services:
# - global
# - dataAnalytic
# - authentication

service=$1
test -z $service && { echo "Usage: $0 service_name"; exit; }

docker build --no-cache --build-arg SERVICE=$service -t cloud_${service,,}:latest  .   
