#!/usr/bin/env bash

# please execute az login before run this script

# dependency: jq - commandline json processor
# https://stedolan.github.io/jq/download/ (Linux, Mac and Windows)
# in windows, execute this script over git bash

resourceGroup='<resource-group-name>'
namespace='<event-hub-name-space>'
eventhub='<event-hub-name>'

sndResult=$(az eventhubs eventhub authorization-rule keys list --resource-group $resourceGroup --namespace-name $namespace --eventhub-name $eventhub --name send)
sndConnStr=$(echo -e $rcvResult | jq '.primaryConnectionString' | xargs)

# send a single message
./ehst -c $sndConnStr -m 'message'

# use -b option to send a single message whose it's content base64 encoded (useful when you capture message with receiver to be replayed)
./ehst -b -c $sndConnStr -m 'bWVzc2FnZQ=='

# send a single message with repetition
./ehst -c $sndConnStr -m 'message' -r 10

# send a single message with repetition and interval(in milliseconds) among the repetitions
./ehst -c $sndConnStr -m 'message' -r 10 -i 500

# send a single message with a property
./ehst -c $sndConnStr -m 'message2' -p "messageId:123123123"

# send a single message with multiple properties
./ehst -c $sndConnStr -m 'message2' -p "messageId:123123123;trackingId:123123123;businessLogicId:123123123;"

# send an amount of messages with batch of messages
./ehst -c $sndConnStr -m 'message' -bm -n 500000

# send a single message with repetition and interval(in milliseconds) among the repetitions
# with batch of messages
./ehst -c $sndConnStr -m 'message' -r 10 -i 500 -bm -n 500000

# replay messages from a file (receiver -c connString > messages-received.json)
./ehst -c $sndConnStr -rp -fp "/mypath/to/the/file/messages-received.json"

# replay messages from a file (receiver -c connString > messages-received.json)
# with repetition and interval among them.
./ehst -c $sndConnStr -r 10 -i 500 -rp -fp "/mypath/to/the/file/messages-received.json"
