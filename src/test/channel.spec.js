var Gusher  = require('../javascripts/gusher'),
    Channel = require('../javascripts/channel');

describe('Channel', function () {

  var channel1;

  beforeEach(function (done) {
    gusher1 = new Gusher(); //stub connection at some point
    gusher1.connection.bind('connected', function () {
      channel1 = gusher1.subscribe('test_channel');
      done();
    });
  });

  describe('constructor', function () {

    it('has a channel name', function () {
      expect(channel1.name).toEqual('test_channel');
    });

  });

});
