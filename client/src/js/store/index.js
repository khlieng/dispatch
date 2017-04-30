import { createStore, applyMiddleware, compose } from 'redux';
import thunk from 'redux-thunk';
import { routerMiddleware } from 'react-router-redux';
import reducer from '../reducers';
import createSocketMiddleware from '../middleware/socket';
import commands from '../commands';

export default function configureStore(socket, history) {
  // eslint-disable-next-line no-underscore-dangle
  const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

  return createStore(reducer, composeEnhancers(
    applyMiddleware(
      routerMiddleware(history),
      thunk,
      createSocketMiddleware(socket),
      commands
    )
  ));
}
