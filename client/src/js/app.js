var React = require('react');

require('./irc');
var socket = require('./socket');
var util = require('./util');
var App = require('./components/App.jsx');
var tabActions = require('./actions/tab');
var serverActions = require('./actions/server');
var channelActions = require('./actions/channel');

var uuid = localStorage.uuid || (localStorage.uuid = util.UUID());
var nick = 'test' + Math.floor(Math.random() * 99999);

socket.on('connect', function() {
	socket.send('uuid', uuid);

	serverActions.connect('irc.freenode.net', nick, 'username', true);
	serverActions.connect('irc.quakenet.org', nick, 'username');

	channelActions.join(['#stuff'], 'irc.freenode.net');
	channelActions.join(['#herp'], 'irc.quakenet.org');
});

socket.on('error', function(error) {
	console.log(error.server + ': ' + error.message);
});

React.render(<App />, document.body);