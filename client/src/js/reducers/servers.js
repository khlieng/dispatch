import { Map, Record } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from '../actions';
import forEach from 'lodash/collection/forEach';

const Server = Record({
  nick: null,
  name: null,
  connected: false
});

export default createReducer(Map(), {
  [actions.CONNECT](state, action) {
    const { host, nick, options } = action;

    if (!state.has(host)) {
      return state.set(host, new Server({
        nick,
        name: options.name || host
      }));
    }

    return state;
  },

  [actions.DISCONNECT](state, action) {
    return state.delete(action.server);
  },

  [actions.SET_NICK](state, action) {
    const { server, nick } = action;
    return state.update(server, s => s.set('nick', nick));
  },

  [actions.SOCKET_SERVERS](state, action) {
    if (!action.data) {
      return state;
    }

    return state.withMutations(s => {
      action.data.forEach(server => {
        s.set(server.host, new Server(server));
      });
    });
  },

  [actions.SOCKET_CONNECTION_UPDATE](state, action) {
    return state.withMutations(s => forEach(action, (connected, server) => {
      if (s.has(server)) {
        s.setIn([server, 'connected'], connected);
      }
    }));
  }
});
