import React from 'react';
import { render } from 'react-dom';
import { browserHistory } from 'react-router';
import configureStore from './store';
import createRoutes from './routes';
import Socket from './util/Socket';
import handleSocket from './socket';
import Root from './containers/Root';

const host = __DEV__ ? `${window.location.hostname}:1337` : window.location.host;
const socket = new Socket(host);

const store = configureStore(socket, browserHistory);
handleSocket(socket, store);

const routes = createRoutes();

render(
  <Root store={store} routes={routes} history={browserHistory} />,
  document.getElementById('root')
);
