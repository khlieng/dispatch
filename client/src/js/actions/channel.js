var Reflux = require('reflux');
var _ = require('lodash');

var socket = require('../socket');

var channelActions = Reflux.createActions([
	'join',
	'part',
	'invite',
	'kick',
	'addUser',
	'removeUser',
	'removeUserAll',
	'renameUser',
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

channelActions.invite.preEmit = function(user, channel, server) {
	socket.send('invite', {
		server: server,
		channel: channel,
		user: user
	});
};

channelActions.kick.preEmit = function(user, channel, server) {
	socket.send('kick', {
		server: server,
		channel: channel,
		user: user
	});
};

module.exports = channelActions;