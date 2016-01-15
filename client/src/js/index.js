import React from 'react';
import { render } from 'react-dom';
import { syncReduxAndRouter } from 'redux-simple-router';
import createBrowserHistory from 'history/lib/createBrowserHistory';
import configureStore from './store';
import createRoutes from './routes';
import Socket from './util/Socket';
import handleSocket from './socket';
import Root from './containers/Root';

const host = __DEV__ ? `${window.location.hostname}:1337` : window.location.host;

const socket = new Socket(host);
const store = configureStore(socket);
handleSocket(socket, store);

const history = createBrowserHistory();
syncReduxAndRouter(history, store);

const routes = createRoutes();

render(<Root store={store} routes={routes} history={history} />, document.getElementById('root'));
