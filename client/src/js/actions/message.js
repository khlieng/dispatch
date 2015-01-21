var Reflux = require('reflux');

var socket = require('../socket.js');

var messageActions = Reflux.createActions([
	'send',
	'add'
]);

messageActions.send.preEmit = function(message, to, server) {
	socket.send('chat', {
		server: server,
		to: to,
		message: message
	});
};

socket.on('message', function(data) {
	messageActions.add(data);
});

socket.on('pm', function(data) {
	messageActions.add(data);
});

socket.on('join', function(data) {
	messageActions.add({
		server: data.server,
		from: '',
		to: data.channels[0],
		message: data.user + ' joined the channel',
		type: 'info'
	});
});

socket.on('part', function(data) {
	messageActions.add({
		server: data.server,
		from: '',
		to: data.channels[0],
		message: data.user + ' left the channel',
		type: 'info'
	});
});

socket.on('motd', function(data) {
	_.each(data.content.split('\n'), function(line) {
		messageActions.add({
			server: data.server,
			from: '',
			to: data.server,
			message: line
		});
	});
});

module.exports = messageActions;