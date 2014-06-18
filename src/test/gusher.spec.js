var Gusher  = require('../javascripts/gusher'),
    Channel = require('../javascripts/channel');

describe('Gusher', function () {

  var gusher1;

  beforeEach(function (done) {
    gusher1 = new Gusher(); //stub connection at some point
    gusher1.connection.bind('connected', function () {
      done();
    });
  });

  describe('#subscribe', function () {

    it('returns a channel', function () {
      var channel = gusher1.subscribe('test_channel');
      expect(channel instanceof Channel).toBe(true);
    });

    it('adds to the list of subscribed channels', function () {
      var channel = gusher1.subscribe('test_channel');
      expect(gusher1.allChannels()).toEqual([channel]);
    });

  });

});
