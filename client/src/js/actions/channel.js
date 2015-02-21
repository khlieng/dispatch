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

channelActions.join.preEmit = (channels, server) => {
	socket.send('join', { server, channels });
};

channelActions.part.preEmit = (channels, server) => {
	socket.send('part', { server, channels });
};

channelActions.invite.preEmit = (user, channel, server) => {
	socket.send('invite', { server, channel, user });
};

channelActions.kick.preEmit = (user, channel, server) => {
	socket.send('kick', { server, channel, user });
};

module.exports = channelActions;