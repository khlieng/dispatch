import React from 'react';
import { render } from 'react-dom';
import { AppContainer } from 'react-hot-loader';
import 'react-virtualized/styles.css';

import configureStore from './store';
import initRouter from './util/router';
import routes from './routes';
import Socket from './util/Socket';
import Root from './containers/Root';
import runModules from './modules';

const host = DEV ? `${window.location.hostname}:1337` : window.location.host;
const socket = new Socket(host);
const store = configureStore(socket);

initRouter(routes, store);
runModules({ store, socket });

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
