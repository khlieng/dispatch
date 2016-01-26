import React from 'react';
import { render } from 'react-dom';
import { browserHistory } from 'react-router';
import { routeActions } from 'redux-simple-router';
import configureStore from './store';
import createRoutes from './routes';
import Socket from './util/Socket';
import handleSocket from './socket';
import Root from './containers/Root';

const host = __DEV__ ? `${window.location.hostname}:1337` : window.location.host;
const socket = new Socket(host);

const store = configureStore(socket, browserHistory);

if (window.__ENV__.servers) {
  store.dispatch({
    type: 'SOCKET_SERVERS',
    data: window.__ENV__.servers
  });
} else {
  store.dispatch(routeActions.replace('/connect'));
}

if (window.__ENV__.channels) {
  store.dispatch({
    type: 'SOCKET_CHANNELS',
    data: window.__ENV__.channels
  });
}

if (window.__ENV__.users) {
  store.dispatch({
    type: 'SOCKET_USERS',
    ...window.__ENV__.users
  });
}

handleSocket(socket, store);

const routes = createRoutes();

render(
  <Root store={store} routes={routes} history={browserHistory} />,
  document.getElementById('root')
);
