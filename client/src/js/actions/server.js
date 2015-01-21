var Reflux = require('reflux');

var socket = require('../socket');

var serverActions = Reflux.createActions([
	'connect',
	'disconnect'
]);

serverActions.connect.preEmit = function(server, nick, username) {
	socket.send('connect', {
		server: server,
		nick: nick,
		username: username
	});
};

serverActions.disconnect.preEmit = function(server) {
	socket.send('quit', { server: server });
};

module.exports = serverActions;