var _            = require('lodash'),
    EventEmitter = require('events').EventEmitter;

var Bus = function () {
};

_.extend(Bus.prototype, EventEmitter.prototype);
Bus.prototype.bind = EventEmitter.prototype.on;

module.exports = Bus;
