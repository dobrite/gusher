var process = require('process');

var dev    = process.env.NODE_ENV !== "production",
    scheme = 'http',
    url    = (dev) ? "localhost" : "example.com",
    port   = (dev) ? 3000 : 80,
    fqd    = scheme + "://" + url + ":" + port.toString();

module.exports = {
  scheme: scheme,
  url: url,
  port: port,
  fqd: fqd
};
