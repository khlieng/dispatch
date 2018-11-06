import React from 'react';
import { render } from 'react-dom';

import Root from 'components/Root';
import { appSet } from 'state/app';
import initRouter from 'utils/router';
import Socket from 'utils/Socket';
import configureStore from './store';
import routes from './routes';
import runModules from './modules';
import { register } from './serviceWorker';
import '../css/fonts.css';
import '../css/fontello.css';
import '../css/style.css';

const production = process.env.NODE_ENV === 'production';
const host = production
  ? window.location.host
  : `${window.location.hostname}:1337`;
const socket = new Socket(host);
const store = configureStore(socket);

initRouter(routes, store);
runModules({ store, socket });

render(<Root store={store} />, document.getElementById('root'));

register({
  onUpdate: () => store.dispatch(appSet('newVersionAvailable', true))
});
