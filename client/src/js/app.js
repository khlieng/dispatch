import React from 'react';
import { render } from 'react-dom';
import { Router, Route, IndexRoute } from 'react-router';
import createBrowserHistory from 'history/lib/createBrowserHistory';
import './irc';
import './command';
import socket from './socket';
import util from './util';
import App from './components/App.jsx';
import Connect from './components/Connect.jsx';
import Chat from './components/Chat.jsx';
import Settings from './components/Settings.jsx';
import routeActions from './actions/route';

let uuid = localStorage.uuid;
if (!uuid) {
	routeActions.navigate('connect', true);
	localStorage.uuid = uuid = util.UUID();
}

socket.on('connect', () => socket.send('uuid', uuid));
socket.on('error', error => console.log(error.server + ': ' + error.message));

const routes = (
	<Route path="/" component={App}>
		<Route path="connect" component={Connect} />
		<Route path="settings" component={Settings} />
		<Route path="/:server" component={Chat} />
		<Route path="/:server/:channel" component={Chat} />
		<IndexRoute component={Settings} />
	</Route>
);
const history = createBrowserHistory();

render(<Router routes={routes} history={history} />, document.getElementById('root'));
