var _            = require('lodash'),
    EventEmitter = require('events').EventEmitter,
    structures   = require('./structures');

var Channel = function (name, bus) {
  this.name = name;
  this.bus = bus;
  var msg = structures.subscribe(name);
  this.bus.send(JSON.stringify(msg));
};

Channel.prototype.bind = function(eventName, callback) {
  this.bus.bind(this.name + ":" + eventName, callback);
};

module.exports = Channel;
