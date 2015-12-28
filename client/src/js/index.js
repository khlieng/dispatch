import React from 'react';
import { render } from 'react-dom';
import { syncReduxAndRouter, replacePath } from 'redux-simple-router';
import createBrowserHistory from 'history/lib/createBrowserHistory';
import configureStore from './store';
import createRoutes from './routes';
import Socket from './util/Socket';
import handleSocket from './socket';
import { createUUID } from './util';
import Root from './containers/Root';

const socket = __DEV__ ?
  new Socket(`${window.location.hostname}:1337`) :
  new Socket(window.location.host);

const store = configureStore(socket);
const routes = createRoutes();
const history = createBrowserHistory();

syncReduxAndRouter(history, store);
handleSocket(socket, store);

let uuid = localStorage.uuid;
if (!uuid) {
  store.dispatch(replacePath('/connect'));
  localStorage.uuid = uuid = createUUID();
}

socket.on('connect', () => socket.send('uuid', uuid));

render(<Root store={store} routes={routes} history={history} />, document.getElementById('root'));
