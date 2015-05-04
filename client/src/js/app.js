var React = require('react');
var Router = require('react-router');
var Route = Router.Route;
var DefaultRoute = Router.DefaultRoute;

require('./irc');
require('./command');
var socket = require('./socket');
var util = require('./util');
var App = require('./components/App.jsx');
var Connect = require('./components/Connect.jsx');
var Chat = require('./components/Chat.jsx');
var Settings = require('./components/Settings.jsx');
var routeActions = require('./actions/route');

var uuid = localStorage.uuid;
if (!uuid) {
	routeActions.navigate('connect', true);
	localStorage.uuid = uuid = util.UUID();
}

socket.on('connect', () => socket.send('uuid', uuid));
socket.on('error', (error) => console.log(error.server + ': ' + error.message));

var routes = (
	<Route name="app" path="/" handler={App}>
		<Route name="connect" handler={Connect} />
		<Route name="settings" handler={Settings} />
		<Route name="status" path="/:server" handler={Chat} />
		<Route name="chat" path="/:server/:channel" handler={Chat} />
		<DefaultRoute handler={Settings} />
	</Route>
);

Router.run(routes, Router.HistoryLocation, (Handler) => {
	React.render(<Handler />, document.body);
});