import { Map, Record } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from '../actions';

const Server = Record({
  nick: null,
  name: null
});

export default createReducer(Map(), {
  [actions.CONNECT](state, action) {
    let { server } = action;
    const { nick, options } = action;

    const i = server.indexOf(':');
    if (i > 0) {
      server = server.slice(0, i);
    }

    return state.set(server, new Server({
      nick,
      name: options.name || server
    }));
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
        s.set(server.address, new Server(server));
      });
    });
  }
});
