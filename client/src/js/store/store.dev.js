import { createStore, applyMiddleware, compose } from 'redux';
import thunk from 'redux-thunk';
import reducer from '../reducers';
import createSocketMiddleware from '../middleware/socket';
import commands from '../commands';
import DevTools from '../containers/DevTools';

export default function configureStore(socket, initialState) {
  const finalCreateStore = compose(
    applyMiddleware(
      thunk,
      createSocketMiddleware(socket),
      commands
    ),
    DevTools.instrument()
  )(createStore);

  const store = finalCreateStore(reducer, initialState);

  if (module.hot) {
    module.hot.accept('../reducers', () => {
      store.replaceReducer(require('../reducers').default);
    });
  }

  return store;
}
