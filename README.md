# message queue websocket

Message queue websocket is a complete websocket application featuring consuming messages from message brokers like Kafka, RabbitMQ, etc and dispatching to end client.

### Table of Contents

1. [Introduction](#introduction)
2. [Running the project](#running)
    1. [Make](#by-make)
    2. [Docker](#by-docker)
3. [Private channel guide](#private-channel-feature)
    1. [Client Side](#client-side)
    2. [Backend side](#backend-side)
4. [Kafka message contract](#kafka-message-contract)

### Introduction

The whole structure is designed in clean architecture. To more information read
this [book](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164). We have tries to
have a wel designed cloud native application as an art of work to encourage other services comply a maintainable,
readable, scalable and cloud friendly one. Any comments to improving and fixing security issues would be appreciated.

### Running

#### By make

Form the root directory you can run this command to run:

```bash
make start
```

#### By docker

Form the root directory you can run this command to run the project in dev mode:

```bash
docker compose -f development/docker-compose.yml up
```

### Private channel feature

#### Client side

To communicate with channel, you must subscribe to the channel by publishing a message with specific format.

```json
{
  "action": "subscribe",
  "channel": "channel-name"
}
```

If you want unsubscribe just change the action to `unsubscribe`.

*Although to authenticate, client should send `token` like other P2P services in the header.

```js
const WebSocket = require('ws');

var socket = new WebSocket("wss://abansite.com/wss/private");

socket.onopen = () => {
    socket.send(JSON.stringify({"action": "subscribe", "channel": "channel-name", "data": "this is optional."}))
    //socket.send(JSON.stringify({"action": "unsubscribe", "channel": "channel-name", "data": "this is optional."}))
}

socket.onmessage = function (event) {
    console.log(event.data)
};
```

#### Backend side

The kafka message format for a specific channel should be in the below format.
`user_id` should be passed into the message's header and service name (which is `channel-name`) should be the key of the
messages.
All `user_id` provided messages would be considered as a *private message*:

```json
header = {
  "user-id": "123414",
  ...
}

key = "channel-name"

value = {
  ...
}
```

### Kafka Message Contract

Sample message from `websocket_public_topic` which we support to process it:

```js
header = {
    "status_code": "None",
    "correlation_id": null,
    "id": "a4595f74-c102-4d59-bccb-e5820440d07b",
    "source": "p2p",
    "gateway_topic": null
}

key = "orderbook_socket"

value = {
    "key": "orderbook",
    "value": {
        "market": "BTC/IRT",
        "bids": [
            {
                "price": "99000000.0",
                "volume": "0.0020000"
            },
            {
                "price": "98000000.0",
                "volume": "0.0800000"
            },
            {
                "price": "10000000.0",
                "volume": "17.3456062"
            }
        ],
        "asks": [
            {
                "price": "110000000.0",
                "volume": "0.9800000"
            },
            {
                "price": "100000000.0",
                "volume": "0.2102898"
            }
        ]
    }
}
```