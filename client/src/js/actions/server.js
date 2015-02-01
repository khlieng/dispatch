var Reflux = require('reflux');

var socket = require('../socket');

var serverActions = Reflux.createActions([
	'connect',
	'disconnect',
	'load'
]);

serverActions.connect.preEmit = function(server, nick, username, tls) {
	socket.send('connect', {
		server: server,
		nick: nick,
		username: username,
		tls: tls || false
	});
};

serverActions.disconnect.preEmit = function(server) {
	socket.send('quit', { server: server });
};

module.exports = serverActions;