var Reflux = require('reflux');

var sock = require('../socket.js')('/ws');

var messageActions = Reflux.createActions([
	'send',
	'add',
	'selectTab'
]);

messageActions.send.preEmit = function(message, to, server) {
	sock.send('chat', {
		server: server,
		to: to,
		message: message
	});
};

module.exports = messageActions;