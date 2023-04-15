# websocket

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