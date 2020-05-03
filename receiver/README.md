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