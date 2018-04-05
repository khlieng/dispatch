import React from 'react';
import { render } from 'react-dom';
import 'react-virtualized/styles.css';

import Root from 'components/Root';
import initRouter from 'utils/router';
import Socket from 'utils/Socket';
import configureStore from './store';
import routes from './routes';
import runModules from './modules';

const production = process.env.NODE_ENV === 'production';
const host = production ? window.location.host : `${window.location.hostname}:1337`;
const socket = new Socket(host);
const store = configureStore(socket);

initRouter(routes, store);
runModules({ store, socket });

render(<Root store={store} />, document.getElementById('root'));
