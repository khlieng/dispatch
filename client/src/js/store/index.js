import { createStore, applyMiddleware, compose } from 'redux';
import thunk from 'redux-thunk';
import { routerMiddleware } from 'react-router-redux';
import reducer from '../reducers';
import createSocketMiddleware from '../middleware/socket';
import commands from '../commands';

export default function configureStore(socket, history, initialState) {
  const finalCreateStore = compose(
    applyMiddleware(
      routerMiddleware(history),
      thunk,
      createSocketMiddleware(socket),
      commands
    ),
    window.devToolsExtension ? window.devToolsExtension() : f => f
  )(createStore);

  return finalCreateStore(reducer, initialState);
}
