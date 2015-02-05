var React = require('react');
var Router = require('react-router');
var Route = Router.Route;
var DefaultRoute = Router.DefaultRoute;

require('./irc');
var socket = require('./socket');
var util = require('./util');
var App = require('./components/App.jsx');
var Connect = require('./components/Connect.jsx');
var Chat = require('./components/Chat.jsx');
var Settings = require('./components/Settings.jsx');

var uuid = localStorage.uuid || (localStorage.uuid = util.UUID());

socket.on('connect', function() {
	socket.send('uuid', uuid);
});

socket.on('error', function(error) {
	console.log(error.server + ': ' + error.message);
});

var routes = (
	<Route name="app" path="/" handler={App}>
		<Route name="connect" handler={Connect} />
		<Route name="settings" handler={Settings} />
		<DefaultRoute handler={Chat} />
	</Route>
);

Router.run(routes, Router.HistoryLocation, function(Handler) {
	React.render(<Handler />, document.body);
});