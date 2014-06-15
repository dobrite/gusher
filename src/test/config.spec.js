describe('config', function () {

  var config;

  beforeEach(function () {
    config = require('../javascripts/config');
  });

  context('given when NODE_ENV is development', function () {
    it('defines scheme', function () {
      expect(config.scheme).to.equal('http');
    });

    it('constructs url according to env', function () {
      expect(config.url).to.equal('localhost');
    });

    it('constructs port according to env', function () {
      expect(config.port).to.equal(3000);
    });

    it('constructs fqd according to env', function () {
      expect(config.fqd).to.equal('http://localhost:3000');
    });

  });

});
