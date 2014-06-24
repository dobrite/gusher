var Promise = require('es6-promise').Promise,
    Gusher  = require('../javascripts/gusher'),
    API     = require('../javascripts/utils');

var setupClients = function(num, done) {
  var count = 0;

  return Array.apply(null, Array(num)).map(function () {
    var gusher = new Gusher();
    gusher.connection.bind('connected', function () {
      if (++count === num) done();
    });
    return gusher;
  });
};

var teardownClients = function(clients, done) {
  var count = 0,
      clientsLength = clients.length;

  return clients.map(function (gusher) {
    gusher.connection.bind('disconnected', function () {
      if (++count === clientsLength) done();
    });
    gusher.disconnect();
  });
};

describe('Integration', function () {

  var clients;

  beforeEach(function (done) {
    clients = setupClients(2, done);
  });

  afterEach(function (done) {
    teardownClients(clients, done);
  });

  it('publishes the message', function (done) {
    var channel = clients[0].subscribe('test-channel');
    channel.bind('test-event', function (data) {
      expect(data.message).toEqual('yo!');
      done();
    });
    API.post('test-channel', 'test-event', { message: "yo!" });
  });

  it('publishes the message to everyone', function (done) {
    var promises = clients.map(function (elem) {
      return new Promise(function(resolve) {
        elem.subscribe('test-channel').bind('test-event', resolve);
      });
    });

    Promise.all(promises).then(function(results) {
      results.map(function (elem) {
        expect(elem.message).toEqual('yo!');
      });
      done();
    });

    API.post('test-channel', 'test-event', { message: "yo!" });
  });

});
