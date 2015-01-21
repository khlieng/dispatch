var Reflux = require('reflux');

var sock = require('../socket.js')('/ws');

var channelActions = Reflux.createActions([
	'join',
	'joined',
	'part',
	'parted',
	'quit',
	'setUsers',
	'setTopic',
	'load'
]);

channelActions.join.preEmit = function(data) {
	sock.send('join', data);
};

channelActions.part.preEmit = function(data) {
	sock.send('part', data);
};

sock.on('join', function(data) {
	channelActions.joined(data.user, data.server, data.channels[0]);
});

sock.on('part', function(data) {
	channelActions.parted(data.user, data.server, data.channels[0]);
});

sock.on('quit', function(data) {
	channelActions.quit(data.user, data.server);
});

sock.on('users', function(data) {
	channelActions.setUsers(data.users, data.server, data.channel);
});

sock.on('topic', function(data) {
	channelActions.setTopic(data.topic, data.server, data.channel);
});

sock.on('channels', function(data) {
	channelActions.load(data);
});

module.exports = channelActions;