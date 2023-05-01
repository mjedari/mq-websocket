# websocket

### Private channel request

#### Client side

To communicate with channel, you must subscribe to the channel by publishing a message with specific format.
```json
{
  "action" : "subscribe",
  "channel" : "channel-name"
}
```
If you want unsubscribe just change the action to `unsubscribe`.

*Although to authenticate, client should send `token` like other P2P services in the header. 

```js
const WebSocket = require('ws');

var socket = new WebSocket("ws://localhost:8000/private");

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
`user_id` should be passed into the message's header and service name (which is `channel-name`) should be the key of the messages.
All `user_id` provided messages would be considered as a *private message*:
```json
header = {
  "user-id" : "123414",
   ...
}

key = "channel-name"

value = {
  ...
}
```

### Sample message from `websocket_public_topic`

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