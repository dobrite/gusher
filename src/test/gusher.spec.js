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

  afterEach(function (done) {
    gusher1.connection.bind('disconnected', function () {
      done();
    });
    gusher1.disconnect();
  });

  describe('constructor', function () {
    it('returns a Gusher', function () {
      expect(gusher1 instanceof Gusher).toBe(true);
    });

    it('is not susbscribed to any channels', function () {
      expect(gusher1.allChannels()).toEqual([]);
    });

    it('stores options', function () {
      gusher1 = new Gusher('', {ping: "pong"});
      expect(gusher1.options).toEqual({ping: "pong"});
    });
  });

  describe('#subscribe', function () {

    it('returns a Channel', function () {
      var channel = gusher1.subscribe('test_channel');
      expect(channel instanceof Channel).toBe(true);
    });

    it('adds to the list of subscribed channels', function () {
      var channel1 = gusher1.subscribe('test_channel');
      var channel2 = gusher1.subscribe('test_channel2');
      expect(gusher1.allChannels()).toEqual([channel1, channel2]);
    });

    it('does not add an already subscribed channel', function () {
      var channel1 = gusher1.subscribe('test_channel');
      var channel2 = gusher1.subscribe('test_channel');
      expect(gusher1.allChannels()).toEqual([channel1]);
    });
  });

  describe('#unsubscribe', function () {

    iit('does not error when called on a channel that is not subscribed', function () {
      gusher1.unsubscribe('test_channel');
      expect(1).toEqual(1);
    });
  });
});
