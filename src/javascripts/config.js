var process = require('process');

var dev = process.env.NODE_ENV !== "production";

var config = {
  url:  (dev) ? "localhost" : "example.com",
  port: (dev) ? 3000 : 80,
  fqd: function () {
    return "http://" + this.url + ":" + this.port.toString() + "/gusher/";
  }()
};

module.exports = config;
