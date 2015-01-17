var Reflux = require('reflux');
var sock = require('../socket.js')('/ws');

var serverActions = Reflux.createActions([
	'connect',
	'disconnect'
]);

serverActions.connect.preEmit = function(data) {
	sock.send('connect', data);
};

module.exports = serverActions;