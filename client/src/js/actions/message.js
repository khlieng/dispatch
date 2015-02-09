var Reflux = require('reflux');
var _ = require('lodash');

var socket = require('../socket');

var messageActions = Reflux.createActions([
	'send',
	'add',
	'broadcast',
	'inform',
	'command'
]);

messageActions.send.preEmit = function(message, to, server) {
	socket.send('chat', {
		server: server,
		to: to,
		message: message
	});
};

module.exports = messageActions;