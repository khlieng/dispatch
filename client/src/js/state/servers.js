import { Map, Record } from 'immutable';
import { createSelector } from 'reselect';
import createReducer from 'util/createReducer';
import { getSelectedTab, updateSelection } from './tab';
import * as actions from './actions';

const Status = Record({
  connected: false,
  error: null
});

const Server = Record({
  nick: '',
  editedNick: null,
  name: '',
  status: new Status()
});

export const getServers = state => state.servers;

export const getCurrentNick = createSelector(
  getServers,
  getSelectedTab,
  (servers, tab) => {
    const editedNick = servers.getIn([tab.server, 'editedNick']);
    if (editedNick === null) {
      return servers.getIn([tab.server, 'nick']);
    }
    return editedNick;
  }
);

export const getCurrentServerName = createSelector(
  getServers,
  getSelectedTab,
  (servers, tab) => servers.getIn([tab.server, 'name'])
);

export const getCurrentServerStatus = createSelector(
  getServers,
  getSelectedTab,
  (servers, tab) => servers.getIn([tab.server, 'status'])
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

  [actions.SET_NICK](state, { server, nick, editing }) {
    if (editing) {
      return state.setIn([server, 'editedNick'], nick);
    } else if (nick === '') {
      return state.setIn([server, 'editedNick'], null);
    }
    return state;
  },

  [actions.socket.NICK](state, { server, oldNick, newNick }) {
    if (!oldNick || oldNick === state.get(server).nick) {
      return state.update(server, s => s
        .set('nick', newNick)
        .set('editedNick', null)
      );
    }
    return state;
  },

  [actions.socket.NICK_FAIL](state, { server }) {
    return state.setIn([server, 'editedNick'], null);
  },

  [actions.socket.SERVERS](state, { data }) {
    if (!data) {
      return state;
    }

    return state.withMutations(s => {
      data.forEach(server => {
        server.status = new Status(server.status);
        s.set(server.host, new Server(server));
      });
    });
  },

  [actions.socket.CONNECTION_UPDATE](state, action) {
    if (state.has(action.server)) {
      return state.setIn([action.server, 'status'], new Status(action));
    }
    return state;
  }
});

export function connect(server, nick, options) {
  let host = server;
  let port;
  const i = server.indexOf(':');
  if (i > 0) {
    host = server.slice(0, i);
    port = server.slice(i + 1);
  }

  return {
    type: actions.CONNECT,
    host,
    nick,
    options,
    socket: {
      type: 'connect',
      data: {
        host,
        port,
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

export function reconnect(server, settings) {
  return {
    type: actions.RECONNECT,
    server,
    settings,
    socket: {
      type: 'reconnect',
      data: {
        ...settings,
        server
      }
    }
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

export function setNick(nick, server, editing) {
  nick = nick.trim().replace(' ', '');

  const action = {
    type: actions.SET_NICK,
    nick,
    server,
    editing
  };

  if (!editing && nick !== '') {
    action.socket = {
      type: 'nick',
      data: {
        newNick: nick,
        server
      }
    };
  }

  return action;
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
