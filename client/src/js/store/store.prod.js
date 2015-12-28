import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import reducer from '../reducers';
import createSocketMiddleware from '../middleware/socket';
import commands from '../commands';

export default function configureStore(socket, initialState) {
  const finalCreateStore = applyMiddleware(
    thunk,
    createSocketMiddleware(socket),
    commands
  )(createStore);

  return finalCreateStore(reducer, initialState);
}
