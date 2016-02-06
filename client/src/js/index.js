import React from 'react';
import { render } from 'react-dom';
import { browserHistory } from 'react-router';
import { routeActions } from 'react-router-redux';
import configureStore from './store';
import createRoutes from './routes';
import Socket from './util/Socket';
import handleSocket from './socket';
import Root from './containers/Root';

import 'react-virtualized/styles.css';

const host = __DEV__ ? `${window.location.hostname}:1337` : window.location.host;
const socket = new Socket(host);

const store = configureStore(socket, browserHistory);

const env = JSON.parse(document.getElementById('env').innerHTML);

// TODO: Handle this properly
window.__ENV__ = {
  defaults: env.defaults
};

if (env.servers) {
  store.dispatch({
    type: 'SOCKET_SERVERS',
    data: env.servers
  });
} else {
  store.dispatch(routeActions.replace('/connect'));
}

if (env.channels) {
  store.dispatch({
    type: 'SOCKET_CHANNELS',
    data: env.channels
  });
}

if (env.users) {
  store.dispatch({
    type: 'SOCKET_USERS',
    ...env.users
  });
}

handleSocket(socket, store);

const routes = createRoutes();

render(
  <Root store={store} routes={routes} history={browserHistory} />,
  document.getElementById('root')
);
