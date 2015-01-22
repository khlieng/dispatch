var Reflux = require('reflux');

var socket = require('../socket');

var serverActions = Reflux.createActions([
	'connect',
	'disconnect',
	'load'
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

socket.on('servers', function(data) {
	serverActions.load(data);
});

module.exports = serverActions;