var structures = {};

structures.subscribe = function(channelName) {
  return {
    "event": "pusher:subscribe",
    "data": JSON.stringify({ //JSON serialized String
      "channel": channelName,
      "auth": "",
      "channel_data": ""
    })
  };
};

structures.unsubscribe = function(channelName) {
  return {
    "event": "pusher:unsubscribe",
    "data": JSON.stringify({ //JSON serialized String
      "channel": channelName,
    })
  };
};

module.exports = structures;
