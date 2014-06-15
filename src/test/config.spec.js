describe('config', function () {

  var config;

  beforeEach(function () {
    config = require('../javascripts/config');
  });

  afterEach(function () {
  });

  context('given when NODE_ENV is development', function () {
    it('constructs url according to env', function () {
      expect(config.url).to.equal('localhost');
    });

    it('constructs port according to env', function () {
      expect(config.port).to.equal(3000);
    });

    it('constructs fqd according to env', function () {
      expect(config.fqd).to.equal('http://localhost:3000/gusher/');
    });

  });

});
