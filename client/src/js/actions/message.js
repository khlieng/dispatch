var Reflux = require('reflux');
var _ = require('lodash');

var socket = require('../socket');

var messageActions = Reflux.createActions([
	'send',
	'add',
	'broadcast'
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

socket.on('quit', function(data) {
	messageActions.broadcast(data.user + ' has quit', data.server);
});

socket.on('nick', function(data) {
	messageActions.broadcast(data.old + ' changed nick to ' + data.new, data.server);
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