var structures = {};

structures.subscribe = function(channelName) {
  return {
    "event": "gusher:subscribe",
    "data": JSON.stringify({ //JSON serialized String
      "channel": channelName,
      "auth": "",
      "channel_data": ""
    })
  };
};

module.exports = structures;
