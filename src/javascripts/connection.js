var _            = require('lodash'),
    EventEmitter = require('events').EventEmitter,
    config       = require('./config');

var onopen = function () {
  this.bus.emit('connection:connected');
};

var onclose = function () {
  this.bus.emit('connection:disconnected');
};

var onmessage = function (data) {
  var message = JSON.parse(data.data);
  var eventName = message.channel + ":" + message.event;
  this.bus.emit(eventName, JSON.parse(message.data));
};

var Connection = function (bus) {
  this.bus = bus;
  this.connection = new SockJS(config.fqd + "/gusher/");
  this.connection.onopen    = onopen.bind(this);
  this.connection.onclose   = onclose.bind(this);
  this.connection.onmessage = onmessage.bind(this);
};

Connection.prototype.bind = function (eventName, callback) {
  this.bus.bind("connection:" + eventName, callback);
};

Connection.prototype.disconnect = function () {
  this.connection.close();
};

module.exports = Connection;
