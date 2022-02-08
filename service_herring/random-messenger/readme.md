# Random Messenger

Meant to be run as a service. This program will delay a random amount of time,
and then send a message to all users using the `wall` command.
The message is randomly picked from a set list.
## Parameters

`-v`: enable verbose

`--message-first`: Send a message before delaying the first time

Examples:
```
go run random-messenger.go -v

go build random-messenger.go && ./random-messenger --message-first
```