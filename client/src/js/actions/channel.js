var Reflux = require('reflux');
var _ = require('lodash');

var socket = require('../socket');

var channelActions = Reflux.createActions([
	'join',
	'part',
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

module.exports = channelActions;