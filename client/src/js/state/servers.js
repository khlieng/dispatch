import { createSelector } from 'reselect';
import get from 'lodash/get';
import createReducer from 'utils/createReducer';
import { getSelectedTab, updateSelection } from './tab';
import * as actions from './actions';

export const getServers = state => state.servers;

export const getCurrentNick = createSelector(
  getServers,
  getSelectedTab,
  (servers, tab) => {
    if (!servers[tab.server]) {
      return;
    }
    const { editedNick } = servers[tab.server];
    if (editedNick === null) {
      return servers[tab.server].nick;
    }
    return editedNick;
  }
);

export const getCurrentServerName = createSelector(
  getServers,
  getSelectedTab,
  (servers, tab) => get(servers, [tab.server, 'name'])
);

export const getCurrentServerStatus = createSelector(
  getServers,
  getSelectedTab,
  (servers, tab) => get(servers, [tab.server, 'status'], {})
);

export default createReducer(
  {},
  {
    [actions.CONNECT](state, { host, nick, name }) {
      if (!state[host]) {
        state[host] = {
          nick,
          editedNick: null,
          name: name || host,
          status: {
            connected: false,
            error: null
          }
        };
      }
    },

    [actions.DISCONNECT](state, { server }) {
      delete state[server];
    },

    [actions.SET_SERVER_NAME](state, { server, name }) {
      state[server].name = name;
    },

    [actions.SET_NICK](state, { server, nick, editing }) {
      if (editing) {
        state[server].editedNick = nick;
      } else if (nick === '') {
        state[server].editedNick = null;
      }
    },

    [actions.socket.NICK](state, { server, oldNick, newNick }) {
      if (!oldNick || oldNick === state[server].nick) {
        state[server].nick = newNick;
        state[server].editedNick = null;
      }
    },

    [actions.socket.NICK_FAIL](state, { server }) {
      state[server].editedNick = null;
    },

    [actions.socket.SERVERS](state, { data }) {
      if (data) {
        data.forEach(({ host, name, nick, status }) => {
          state[host] = { name, nick, status, editedNick: null };
        });
      }
    },

    [actions.socket.CONNECTION_UPDATE](state, { server, connected, error }) {
      if (state[server]) {
        state[server].status.connected = connected;
        state[server].status.error = error;
      }
    }
  }
);

export function connect(config) {
  return {
    type: actions.CONNECT,
    ...config,
    socket: {
      type: 'connect',
      data: config
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
