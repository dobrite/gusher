{
  "event": "pusher:connection_established",
  "data": { //JSON serializied String
    "socket_id": String,
    "activity_timeout": Number
  }
}

{
  "event": "pusher:subscribe",
  "data": { //JSON serialized String
    "channel": String,
    "auth": String,
    "channel_data": Object
  }
}
