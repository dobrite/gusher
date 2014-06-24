var Gusher  = require('../javascripts/gusher');

var setupClients = function(num, done) {
  var count = 0;

  return Array.apply(null, Array(num)).map(function () {
    var gusher = new Gusher();
    gusher.connection.bind('connected', function () {
      if (++count === num) done();
    });
    return gusher;
  });
};

var teardownClients = function(clients, done) {
  var count = 0,
      clientsLength = clients.length;

  return clients.map(function (gusher) {
    gusher.connection.bind('disconnected', function () {
      if (++count === clientsLength) done();
    });
    gusher.disconnect();
  });
};

module.exports = {
  setupClients: setupClients,
  teardownClients: teardownClients
};
