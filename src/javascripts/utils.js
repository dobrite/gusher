var _           = require('lodash'),
    config      = require('./config'),
    querystring = require('querystring');

var PostData = function (channel, event, data) {
  this.channel  = channel;
  this.event    = event;
  this.data     = JSON.stringify(data);
};

var post = function (channel, event, data) {
  var request = new XMLHttpRequest();
  request.open('POST', config.fqd + '/post/', true);
  request.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
  request.send(querystring.stringify(new PostData(channel, event, data)));
};

module.exports = {
  post: post
};
