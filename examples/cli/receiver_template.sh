#!/usr/bin/env bash

# please execute az login before run this script

# dependency: jq - commandline json processor
# https://stedolan.github.io/jq/download/ (Linux, Mac and Windows)
# in windows, execute this script over git bash

resourceGroup='<resource-group-name>'
namespace='<event-hub-name-space>'
eventhub='<event-hub-name>'

#optional: internally is used $Default consumer group, but is recommended use a specific one
# to avoid disconnect listeners from the default consumer group by accident.
consumerGroup='<consumer-group-name>'

rcvResult=$(az eventhubs eventhub authorization-rule keys list --resource-group $resourceGroup --namespace-name $namespace --eventhub-name $eventhub --name listen)
rcvConnStr=$(echo -e $rcvResult | jq '.primaryConnectionString' | xargs)

# using the default consumer group
./ehr -c $rcvConnStr

# specifying a consumer group
# ./ehr -c $rcvConnStr -g $consumerGroup