import { Map, Record } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from '../actions';

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

  [actions.SOCKET_NICK](state, action) {
    const { server, old } = action;
    if (!old || old === state.get(server).nick) {
      return state.update(server, s => s.set('nick', action.new));
    }
    return state;
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
    return state.withMutations(s =>
      Object.keys(action).forEach(server => {
        if (s.has(server)) {
          s.setIn([server, 'connected'], action[server]);
        }
      })
    );
  }
});
