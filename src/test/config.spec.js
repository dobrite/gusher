describe('config', function () {

  var config;

  beforeEach(function () {
    config = require('../javascripts/config');
  });

  describe('given when NODE_ENV is development', function () {
    it('defines scheme', function () {
      expect(config.scheme).toEqual('http');
    });

    it('constructs url according to env', function () {
      expect(config.url).toEqual('localhost');
    });

    it('constructs port according to env', function () {
      expect(config.port).toEqual(3000);
    });

    it('constructs fqd according to env', function () {
      expect(config.fqd).toEqual('http://localhost:3000');
    });

  });

});
