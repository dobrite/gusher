var Gusher    = require('../javascripts/gusher'),
    Channel   = require('../javascripts/channel'),
    testUtils = require('./test-utils');

describe('Channel', function () {

  var client, channel1;

  beforeEach(function (done) {
    client = testUtils.setupClients(1, done);
  });

  afterEach(function (done) {
    testUtils.teardownClients(client, done);
  });

  describe('constructor', function () {
    beforeEach(function (done) {
      channel1 = client[0].subscribe('test_channel');
      done();
    });

    it('has a channel name', function () {
      expect(channel1.name).toEqual('test_channel');
    });
  });
});
