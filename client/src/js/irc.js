var _ = require('lodash');

var socket = require('./socket');
var channelActions = require('./actions/channel');
var messageActions = require('./actions/message');
var serverActions = require('./actions/server');

socket.on('join', function(data) {
	channelActions.addUser(data.user, data.server, data.channels[0]);
	messageActions.inform(data.user + ' joined the channel', data.server, data.channels[0]);
});

socket.on('part', function(data) {
	channelActions.removeUser(data.user, data.server, data.channels[0]);
	messageActions.inform(withReason(data.user + ' left the channel', data.reason), data.server, data.channels[0]);
});

socket.on('quit', function(data) {
	messageActions.broadcast(withReason(data.user + ' quit', data.reason), data.server, data.user);
	channelActions.removeUserAll(data.user, data.server);
});

socket.on('nick', function(data) {
	messageActions.broadcast(data.old + ' changed nick to ' + data.new, data.server, data.old);
	channelActions.renameUser(data.old, data.new, data.server);
});

socket.on('message', function(data) {
	messageActions.add(data);
});

socket.on('pm', function(data) {
	messageActions.add(data);
});

socket.on('motd', function(data) {
	_.each(data.content.split('\n'), function(line) {
		messageActions.add({
			server: data.server,
			to: data.server,
			message: line
		});
	});
});

socket.on('users', function(data) {
	channelActions.setUsers(data.users, data.server, data.channel);
});

socket.on('topic', function(data) {
	channelActions.setTopic(data.topic, data.server, data.channel);
});

socket.on('mode', function(data) {
	channelActions.setMode(data);
});

socket.on('servers', function(data) {
	serverActions.load(data);
});

socket.on('channels', function(data) {
	channelActions.load(data);
});

function withReason(message, reason) {
	return message + (reason ? ' (' + reason + ')' : '');
}