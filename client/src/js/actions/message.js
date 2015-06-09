var Reflux = require('reflux');

var socket = require('../socket');

var messageActions = Reflux.createActions([
	'send',
	'add',
	'addAll',
	'broadcast',
	'inform',
	'command',
	'setWrapWidth'
]);

messageActions.send.preEmit = (message, to, server) => {
	socket.send('chat', { server, to, message });
};

module.exports = messageActions;