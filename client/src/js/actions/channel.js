var Reflux = require('reflux');
var _ = require('lodash');

var socket = require('../socket');

var channelActions = Reflux.createActions([
	'join',
	'part',
	'addUser',
	'removeUser',
	'removeUserAll',
	'setUsers',
	'setTopic',
	'setMode',
	'load'
]);

channelActions.join.preEmit = function(channels, server) {
	socket.send('join', {
		server: server,
		channels: channels
	});
};

channelActions.part.preEmit = function(channels, server) {
	socket.send('part', {
		server: server,
		channels: channels
	});
};

socket.on('join', function(data) {
	channelActions.addUser(data.user, data.server, data.channels[0]);
});

socket.on('part', function(data) {
	channelActions.removeUser(data.user, data.server, data.channels[0]);
});

socket.on('quit', function(data) {
	channelActions.removeUserAll(data.user, data.server);
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

socket.on('channels', function(data) {
	channelActions.load(data);
});

module.exports = channelActions;