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

const host = __DEV__ ? `${window.location.hostname}:1337` : window.location.host;

let uuid = localStorage.uuid;
let newUser = false;
if (!uuid) {
  uuid = createUUID();
  newUser = true;
}

const socket = new Socket(host, uuid);
const store = configureStore(socket);
handleSocket(socket, store);

const history = createBrowserHistory();
syncReduxAndRouter(history, store);

if (newUser) {
  store.dispatch(replacePath('/connect'));
  localStorage.uuid = uuid;
}

const routes = createRoutes();

render(<Root store={store} routes={routes} history={history} />, document.getElementById('root'));
