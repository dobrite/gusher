describe('A test suite', function () {

  var sock;

  beforeEach(function (done) {
    sock = new SockJS('http://localhost:3000/echo/');
    sock.onopen = function (data) {
      done();
    };
  });

  afterEach(function (done) {
    sock.onclose = function () {
      done();
    };
    sock.close();
  });

  it('yo!', function (done) {
    sock.onmessage = function (data) {
      expect(data.type).to.be.equal('message');
      expect(data.data).to.be.equal('yo!');
      done();
    };
    sock.send('yo!');
  });

});
