var Promise   = require('es6-promise').Promise,
    Gusher    = require('../javascripts/gusher'),
    API       = require('../javascripts/utils'),
    testUtils = require('./test-utils');

describe('Integration', function () {

  var clients;

  beforeEach(function (done) {
    clients = testUtils.setupClients(1, done);
  });

  afterEach(function (done) {
    testUtils.teardownClients(clients, done);
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
