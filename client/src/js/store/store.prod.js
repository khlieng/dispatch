import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import { syncHistory } from 'react-router-redux';
import reducer from '../reducers';
import createSocketMiddleware from '../middleware/socket';
import commands from '../commands';

export default function configureStore(socket, history, initialState) {
  const finalCreateStore = applyMiddleware(
    syncHistory(history),
    thunk,
    createSocketMiddleware(socket),
    commands
  )(createStore);

  return finalCreateStore(reducer, initialState);
}
