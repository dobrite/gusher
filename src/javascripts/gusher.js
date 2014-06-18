var _            = require('lodash'),
    EventEmitter = require('events').EventEmitter,
    config       = require('./config'),
    Channel      = require('./channel'),
    Connection   = require('./connection'),
    Bus          = require('./bus');

var subscribedChannels = {};

function Gusher (applicationKey, options) {
  options || (options = {});
  this.options = options;
  this.bus = new Bus();
  this.connection = new Connection(this.bus);
}

Gusher.prototype.subscribe = function (channelName) {
  var channel = new Channel(channelName, this.bus);
  subscribedChannels[channelName] = channel;
  return channel;
};

Gusher.prototype.unsubscribe = function (channelName) {
};

Gusher.prototype.allChannels = function () {
  return _.values(subscribedChannels);
};

Gusher.prototype.disconnect = function () {
  this.connection.disconnect();
};

module.exports = Gusher;
