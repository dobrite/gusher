describe('config', function () {

  var config;

  beforeEach(function (done) {
    config = require('config');
  });

  afterEach(function (done) {
  });

  context('given when NODE_ENV is development', function () {
    it('constructs url according to env', function () {
      expect(config.url).to.equal('localhost');
    });

    it('constructs port according to env', function () {
      expect(config.url).to.equal(3000);
    });

  });

});
