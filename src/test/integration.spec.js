var Gusher = require('../javascripts/gusher'),
    API    = require('../javascripts/utils');

describe('Integration', function () {

  var gusher1, gusher2;

  beforeEach(function (done) {
    gusher1 = new Gusher();
    gusher1.connection.bind('connected', function () {
      gusher2 = new Gusher();
      gusher2.connection.bind('connected', function () {
        done();
      });
    });
  });

  afterEach(function (done) {
    gusher1.connection.bind('disconnected', function () {
      gusher2.connection.bind('disconnected', function () {
        done();
      });
      gusher2.disconnect();
    });
    gusher1.disconnect();
  });

  it('publishes the message', function (done) {
    var channel = gusher1.subscribe('test-channel');
    channel.bind('test-event', function (data) {
      expect(data.message).to.be.equal('yo!');
      done();
    });
    API.post('test-channel', 'test-event', { message: "yo!" });
  });

  it('publishes the message to others', function (done) {
    var channel = gusher2.subscribe('test-channel');
    channel.bind('test-event', function (data) {
      expect(data.message).to.be.equal('yo!');
      done();
    });
    API.post('test-channel', 'test-event', { message: "yo!" });
  });

});
