# Introduction

# Build
![Event-Hubs-Cli](https://github.com/azure-open-tools/event-hubs-cli/workflows/Event-Hubs-Cli/badge.svg)

EVENTHUB_SEND_CONNSTR
EVENTHUB_LISTEN_CONNSTR

to build it locally you can do:
```shell script
make build-local-sender && make build-local-receiver
```

or just build by your own like:
```shell script
#sender
go build -ldflags "-s -w" -o sender/bin/ehst sender/main.go

#receiver
go build -ldflags "-s -w" -o receiver/bin/erst receiver/main.go
```

# Event Hub, event json structure

When you send a "message" to the event hub, you are sending an ```event``` and this ```event``` is a json like this:

```json
{
    "Data": "my message comes here, can be a json{} or a byte array",
    "PartitionKey": null,
    "ID": "MyId",
    "SystemProperties": {
        "SequenceNumber": 1,
        "EnqueuedTime": "2020-02-13T12:54:57.642Z",
        "Offset": 21479629240,
        "PartitionID": null,
        "PartitionKey": null
    }
}
```

Properties is an optional field, and you can also add to enrich your event with metadata like this:

```json
{
    "Data": "my message comes here, can be a json{} or a byte array",
    "PartitionKey": null,
    "Properties": {
        "property1": "15",
        "property2": "15",
        "property3": "15",
        ...
        ...
        ...
    },
    "ID": "MyId",
    "SystemProperties": {
        "SequenceNumber": 1,
        "EnqueuedTime": "2020-02-13T12:54:57.642Z",
        "Offset": 21479629240,
        "PartitionID": null,
        "PartitionKey": null
    }
}
```

You can use the sender (producer) command line to send messages (events) to the event hub with properties as well.

If you don't use Azure Events Hub yet, you can read more about it: [Event Hubs About](https://docs.microsoft.com/en-us/azure/event-hubs/event-hubs-about) and
[Event Hubs Features](https://docs.microsoft.com/en-us/azure/event-hubs/event-hubs-features).

# Use cases

# Sender (producer)

# Receiver (consumer)

You can use this tool to start a listener (consumer) against Azure Event Hubs.

Straight to the point:

```shell script
ehrt -c "<conn-string>" 
```

The command above uses the $Default consumer group. For development environment maybe will not be an issue.
when the listener start, it will print the Data field from the event to the stdout.

In environments where you have many listeners, better if you use a specific consumer group to use this tool as debug 
to avoid disconnect others, you can use the command with -g flag to use your consumer group for debugging:
```shell script
ehrt -g "debug-consumer" -c "<conn-string>" 
```

with a -v flag you can print out the whole Event struct (json)
```shell script
ehrt -v -g "debug-consumer" -c "<conn-string>" 
```
If you desire to avoid print all incoming events, you can filter by the data field or properties field like this:

data filter:
```shell script
ehrt -v -d "my message content filter" -g "debug-consumer" -c "<conn-string>"
```
property filter:
```shell script
ehrt -v -p "property filter" -g "debug-consumer" -c "<conn-string>"
```

You can also provide more than one filter, just providing more than one -d or -p flag, like this:
```shell script
ehrt -v -p "property filter1" -p "property filter2" -p "property filter3" -g "debug-consumer" -c "<conn-string>"
```
you can save all the output sending them to a file like this:
```shell script
ehrt -v -g "<consumer-group>" -c "<conn-string>" > output.json
```
It will not generate a well formatted json, but each line of that file will be a valid json since you are using -v flag.
With that file you can play alongside the sender (producer) to replay all the file content to an event hub you wish.

Write events to multiple json files using -s flag:
```shell script
ehrt -v -s -f ${PWD}/ -p "<filter1>" -g "<consumer-group>" -c "<conn-string>"
```
The command above will generate a json file for each event arrived that match the property filter
