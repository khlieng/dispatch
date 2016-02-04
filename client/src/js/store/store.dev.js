import { createStore, applyMiddleware, compose } from 'redux';
import thunk from 'redux-thunk';
import { syncHistory } from 'react-router-redux';
import reducer from '../reducers';
import createSocketMiddleware from '../middleware/socket';
import commands from '../commands';
import DevTools from '../containers/DevTools';

export default function configureStore(socket, history, initialState) {
  const reduxRouterMiddleware = syncHistory(history);

  const finalCreateStore = compose(
    applyMiddleware(
      reduxRouterMiddleware,
      thunk,
      createSocketMiddleware(socket),
      commands
    ),
    DevTools.instrument()
  )(createStore);

  const store = finalCreateStore(reducer, initialState);

  reduxRouterMiddleware.listenForReplays(store);

  if (module.hot) {
    module.hot.accept('../reducers', () => {
      store.replaceReducer(require('../reducers').default);
    });
  }

  return store;
}
