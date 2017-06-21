import { createStore, applyMiddleware, compose } from 'redux';
import thunk from 'redux-thunk';
import createReducer from 'state';
import { routeReducer, routeMiddleware } from 'util/router';
import message from './middleware/message';
import createSocketMiddleware from './middleware/socket';
import commands from './commands';

export default function configureStore(socket) {
  // eslint-disable-next-line no-underscore-dangle
  const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

  const reducer = createReducer(routeReducer);

  const store = createStore(reducer, composeEnhancers(
    applyMiddleware(
      thunk,
      routeMiddleware,
      createSocketMiddleware(socket),
      message,
      commands
    )
  ));

  return store;
}
