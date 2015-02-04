var Reflux = require('reflux');

var socket = require('../socket');

var serverActions = Reflux.createActions([
	'connect',
	'disconnect',
	'setNick',
	'load'
]);

serverActions.connect.preEmit = function(server, nick, username, tls, name) {
	socket.send('connect', {
		server: server,
		nick: nick,
		username: username,
		tls: tls || false,
		name: name || server
	});
};

serverActions.disconnect.preEmit = function(server) {
	socket.send('quit', { server: server });
};

serverActions.setNick.preEmit = function(nick, server) {
	socket.send('nick', {
		server: server,
		new: nick
	});
};

module.exports = serverActions;