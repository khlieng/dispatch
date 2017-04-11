import React from 'react';
import { render } from 'react-dom';
import { browserHistory } from 'react-router';
import { syncHistoryWithStore, replace } from 'react-router-redux';
import { AppContainer } from 'react-hot-loader';
import 'react-virtualized/styles.css';

import configureStore from './store';
import createRoutes from './routes';
import Socket from './util/Socket';
import handleSocket from './socket';
import Root from './containers/Root';

const host = DEV ? `${window.location.hostname}:1337` : window.location.host;
const socket = new Socket(host);

const store = configureStore(socket, browserHistory);

const env = JSON.parse(document.getElementById('env').innerHTML);

// TODO: Handle this properly
window.ENV = {
  defaults: env.defaults
};

if (env.servers) {
  store.dispatch({
    type: 'SOCKET_SERVERS',
    data: env.servers
  });
} else {
  store.dispatch(replace('/connect'));
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
const history = syncHistoryWithStore(browserHistory, store);

const renderRoot = () => {
  render(
    <AppContainer>
      <Root store={store} routes={routes} history={history} />
    </AppContainer>,
    document.getElementById('root')
  );
};

renderRoot();

if (module.hot) {
  module.hot.accept('./routes', () => renderRoot());
}
