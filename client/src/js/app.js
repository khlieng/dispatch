var React = require('react');
var _ = require('lodash');

var sock = require('./socket')('/ws');
var util = require('./util');
var App = require('./components/App.jsx');
var messageActions = require('./actions/message.js');
var tabActions = require('./actions/tab.js');
var serverActions = require('./actions/server.js');
var channelActions = require('./actions/channel.js');

React.render(<App />, document.body);

var uuid = localStorage.uuid || (localStorage.uuid = util.UUID());

tabActions.select('irc.freenode.net');

sock.on('connect', function() {
	sock.send('uuid', uuid);

	serverActions.connect({
		server: 'irc.freenode.net',
		nick: 'test' + Math.floor(Math.random() * 99999),
		username: 'user'
	});

	channelActions.join({
		server: 'irc.freenode.net',
		channels: [ '#stuff', '#go-nuts' ]
	});
});

channelActions.joined.listen(function(user, server, channel) {
	messageActions.add({
		server: server,
		from: '',
		to: channel,
		message: user + ' joined the channel'
	});
});

channelActions.parted.listen(function(user, server, channel) {
	messageActions.add({
		server: server,
		from: '',
		to: channel,
		message: user + ' left the channel'
	});
});

sock.on('message', function(data) {
	messageActions.add(data);
});

sock.on('pm', function(data) {
	messageActions.add(data);
});

sock.on('topic', function(data) {
	messageActions.add({
		server: data.server,
		from: '',
		to: data.channel,
		message: data.topic
	});
});

sock.on('motd', function(data) {
	_.each(data.content.split('\n'), function(line) {
		messageActions.add({
			server: data.server,
			from: '',
			to: data.server,
			message: line
		});
	});
});