import { createStore, applyMiddleware, compose } from 'redux';
import thunk from 'redux-thunk';
import createReducer from '../reducers';
import { routeReducer, routeMiddleware } from '../util/router';
import createSocketMiddleware from '../middleware/socket';
import commands from '../commands';

export default function configureStore(socket) {
  // eslint-disable-next-line no-underscore-dangle
  const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

  const reducer = createReducer(routeReducer);

  const store = createStore(reducer, composeEnhancers(
    applyMiddleware(
      routeMiddleware,
      thunk,
      createSocketMiddleware(socket),
      commands
    )
  ));

  return store;
}
