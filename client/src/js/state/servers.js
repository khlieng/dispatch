import { Map, Record } from 'immutable';
import { createSelector } from 'reselect';
import createReducer from '../util/createReducer';
import { getSelectedTab, updateSelection } from './tab';
import * as actions from './actions';

const Server = Record({
  nick: null,
  name: null,
  connected: false
});

export const getServers = state => state.servers;

export const getCurrentNick = createSelector(
  getServers,
  getSelectedTab,
  (servers, tab) => servers.getIn([tab.server, 'nick'], '')
);

export const getCurrentServerName = createSelector(
  getServers,
  getSelectedTab,
  (servers, tab) => servers.getIn([tab.server, 'name'], '')
);

export default createReducer(Map(), {
  [actions.CONNECT](state, { host, nick, options }) {
    if (!state.has(host)) {
      return state.set(host, new Server({
        nick,
        name: options.name || host
      }));
    }

    return state;
  },

  [actions.DISCONNECT](state, { server }) {
    return state.delete(server);
  },

  [actions.SET_SERVER_NAME](state, { server, name }) {
    return state.setIn([server, 'name'], name);
  },

  [actions.socket.NICK](state, action) {
    const { server, old } = action;
    if (!old || old === state.get(server).nick) {
      return state.update(server, s => s.set('nick', action.new));
    }
    return state;
  },

  [actions.socket.SERVERS](state, { data }) {
    if (!data) {
      return state;
    }

    return state.withMutations(s => {
      data.forEach(server => {
        s.set(server.host, new Server(server));
      });
    });
  },

  [actions.socket.CONNECTION_UPDATE](state, action) {
    return state.withMutations(s =>
      Object.keys(action).forEach(server => {
        if (s.has(server)) {
          s.setIn([server, 'connected'], action[server]);
        }
      })
    );
  }
});

export function connect(server, nick, options) {
  let host = server;
  const i = server.indexOf(':');
  if (i > 0) {
    host = server.slice(0, i);
  }

  return {
    type: actions.CONNECT,
    host,
    nick,
    options,
    socket: {
      type: 'connect',
      data: {
        server,
        nick,
        username: options.username || nick,
        password: options.password,
        realname: options.realname || nick,
        tls: options.tls || false,
        name: options.name || server
      }
    }
  };
}

export function disconnect(server) {
  return dispatch => {
    dispatch({
      type: actions.DISCONNECT,
      server,
      socket: {
        type: 'quit',
        data: { server }
      }
    });
    dispatch(updateSelection());
  };
}

export function whois(user, server) {
  return {
    type: actions.WHOIS,
    user,
    server,
    socket: {
      type: 'whois',
      data: { user, server }
    }
  };
}

export function away(message, server) {
  return {
    type: actions.AWAY,
    message,
    server,
    socket: {
      type: 'away',
      data: { message, server }
    }
  };
}

export function setNick(nick, server) {
  return {
    type: actions.SET_NICK,
    nick,
    server,
    socket: {
      type: 'nick',
      data: {
        new: nick,
        server
      }
    }
  };
}

export function isValidServerName(name) {
  return name.trim() !== '';
}

export function setServerName(name, server) {
  const action = {
    type: actions.SET_SERVER_NAME,
    name,
    server
  };

  if (isValidServerName(name)) {
    action.socket = {
      type: 'set_server_name',
      data: {
        name,
        server
      },
      debounce: {
        delay: 1000,
        key: server
      }
    };
  }

  return action;
}
