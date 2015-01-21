var React = require('react');

var socket = require('./socket');
var util = require('./util');
var App = require('./components/App.jsx');
var tabActions = require('./actions/tab.js');
var serverActions = require('./actions/server.js');
var channelActions = require('./actions/channel.js');

var uuid = localStorage.uuid || (localStorage.uuid = util.UUID());
var nick = 'test' + Math.floor(Math.random() * 99999);

socket.on('connect', function() {
	socket.send('uuid', uuid);

	serverActions.connect('irc.freenode.net', nick, 'username');
	channelActions.join(['#stuff'], 'irc.freenode.net');
	tabActions.select('irc.freenode.net');
});

React.render(<App />, document.body);