var Gusher = require('../javascripts/gusher'),
    post   = require('../javascripts/utils').post;

describe('A test suite', function () {

  var client_1, client_2;

  beforeEach(function (done) {
    client_1 = new Gusher();
    client_1.connection.bind('connected', function () {
      client_2 = new Gusher();
      client_2.connection.bind('connected', function () {
        done();
      });
    });
  });

  afterEach(function (done) {
    client_1.connection.bind('disconnected', function () {
      client_2.connection.bind('disconnected', function () {
        done();
      });
      client_2.disconnect();
    });
    client_1.disconnect();
  });

  it('echos the message', function (done) {
    var channel = client_1.subscribe('test-channel');
    channel.bind('test-event', function (data) {
      expect(data.message).to.be.equal('yo!');
      done();
    });
    post('test-channel', 'test-event', { message: "yo!" });
  });

  it('publishes the message to others', function (done) {
    var channel = client_2.subscribe('test-channel');
    channel.bind('test-event', function (data) {
      expect(data.message).to.be.equal('yo!');
      done();
    });
    post('test-channel', 'test-event', { message: "yo!" });
  });

});
