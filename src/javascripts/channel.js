var _            = require('lodash'),
    EventEmitter = require('events').EventEmitter;

var Channel = function (name, bus) {
  this.name = name;
  this.bus = bus;
};

Channel.prototype.bind = function(eventName, callback) {
  this.bus.bind(this.name + ":" + eventName, callback);
};

module.exports = Channel;
