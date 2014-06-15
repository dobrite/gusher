var _            = require('lodash'),
    EventEmitter = require('events').EventEmitter;

var Channel = function (channelName, bus) {
  this.channelName = channelName;
  this.bus = bus;
};

Channel.prototype.bind = function(eventName, callback) {
  this.bus.bind(this.channelName + ":" + eventName, callback);
};

module.exports = Channel;
