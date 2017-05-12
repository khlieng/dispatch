import React from 'react';
import { render } from 'react-dom';
import { AppContainer } from 'react-hot-loader';
import 'react-virtualized/styles.css';

import configureStore from './store';
import initRouter, { replace } from './util/router';
import routes from './routes';
import Socket from './util/Socket';
import handleSocket from './socket';
import Root from './containers/Root';
import { addMessages } from './actions/message';
import { initWidthUpdates } from './util/messageHeight';

const host = DEV ? `${window.location.hostname}:1337` : window.location.host;
const socket = new Socket(host);

const store = configureStore(socket);
initRouter(routes, store);
handleSocket(socket, store);

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

initWidthUpdates(store, () => {
  if (env.messages) {
    const { messages, server, to, next } = env.messages;
    store.dispatch(addMessages(messages, server, to, false, next));
  }
});

const renderRoot = () => {
  render(
    <AppContainer>
      <Root store={store} />
    </AppContainer>,
    document.getElementById('root')
  );
};

renderRoot();

if (module.hot) {
  module.hot.accept('./containers/Root', () => renderRoot());
}
