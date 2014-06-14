describe('A test suite', function () {

  var sock1, sock2;

  beforeEach(function (done) {
    sock1 = new SockJS('http://localhost:3000/gusher/');
    sock1.onopen = function () {
      sock2 = new SockJS('http://localhost:3000/gusher/');
      sock2.onopen = function () {
        done();
      };
    };
  });

  afterEach(function (done) {
    sock1.onclose = function () {
      sock2.onclose = function () {
        done();
      };
      sock2.close();
    };
    sock1.close();
  });

  it('echos the message', function (done) {
    sock1.onmessage = function (data) {
      expect(data.type).to.be.equal('message');
      expect(data.data).to.be.equal('yo!');
      done();
    };
    sock1.send('yo!');
  });

  it('publishes the message to others', function (done) {
    sock2.onmessage = function (data) {
      expect(data.type).to.be.equal('message');
      expect(data.data).to.be.equal('yo!');
      done();
    };
    sock1.send('yo!');
  });

});
