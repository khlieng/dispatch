import { createSelector } from 'reselect';
import get from 'lodash/get';
import createReducer from 'utils/createReducer';
import { getSelectedTab, updateSelection } from './tab';
import * as actions from './actions';

export const getNetworks = state => state.networks;

export const getCurrentNick = createSelector(
  getNetworks,
  getSelectedTab,
  (networks, tab) => {
    if (!networks[tab.network]) {
      return;
    }
    const { editedNick } = networks[tab.network];
    if (editedNick === null) {
      return networks[tab.network].nick;
    }
    return editedNick;
  }
);

export const getCurrentNetworkName = createSelector(
  getNetworks,
  getSelectedTab,
  (networks, tab) => get(networks, [tab.network, 'name'])
);

export const getCurrentNetworkError = createSelector(
  getNetworks,
  getSelectedTab,
  (networks, tab) => get(networks, [tab.network, 'error'], null)
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
          connected: false,
          error: null,
          features: {}
        };
      }
    },

    [actions.DISCONNECT](state, { network }) {
      delete state[network];
    },

    [actions.SET_NETWORK_NAME](state, { network, name }) {
      state[network].name = name;
    },

    [actions.SET_NICK](state, { network, nick, editing }) {
      if (editing) {
        state[network].editedNick = nick;
      } else if (nick === '') {
        state[network].editedNick = null;
      }
    },

    [actions.socket.NICK](state, { network, oldNick, newNick }) {
      if (!oldNick || oldNick === state[network].nick) {
        state[network].nick = newNick;
        state[network].editedNick = null;
      }
    },

    [actions.socket.NICK_FAIL](state, { network }) {
      state[network].editedNick = null;
    },

    [actions.INIT](state, { networks }) {
      if (networks) {
        networks.forEach(
          ({ host, name = host, nick, connected, error, features = {} }) => {
            state[host] = {
              name,
              nick,
              connected,
              error,
              features,
              editedNick: null
            };
          }
        );
      }
    },

    [actions.socket.CONNECTION_UPDATE](state, { network, connected, error }) {
      if (state[network]) {
        state[network].connected = connected;
        state[network].error = error;
      }
    },

    [actions.socket.FEATURES](state, { network, features }) {
      const srv = state[network];
      if (srv) {
        srv.features = features;

        if (features.NETWORK && srv.name === network) {
          srv.name = features.NETWORK;
        }
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

export function disconnect(network) {
  return dispatch => {
    dispatch({
      type: actions.DISCONNECT,
      network,
      socket: {
        type: 'quit',
        data: { network }
      }
    });
    dispatch(updateSelection());
  };
}

export function reconnect(network, settings) {
  return {
    type: actions.RECONNECT,
    network,
    settings,
    socket: {
      type: 'reconnect',
      data: {
        ...settings,
        network
      }
    }
  };
}

export function whois(user, network) {
  return {
    type: actions.WHOIS,
    user,
    network,
    socket: {
      type: 'whois',
      data: { user, network }
    }
  };
}

export function away(message, network) {
  return {
    type: actions.AWAY,
    message,
    network,
    socket: {
      type: 'away',
      data: { message, network }
    }
  };
}

export function setNick(nick, network, editing) {
  nick = nick.trim().replace(' ', '');

  const action = {
    type: actions.SET_NICK,
    nick,
    network,
    editing
  };

  if (!editing && nick !== '') {
    action.socket = {
      type: 'nick',
      data: {
        newNick: nick,
        network
      }
    };
  }

  return action;
}

export function isValidNetworkName(name) {
  return name.trim() !== '';
}

export function setNetworkName(name, network) {
  const action = {
    type: actions.SET_NETWORK_NAME,
    name,
    network
  };

  if (isValidNetworkName(name)) {
    action.socket = {
      type: 'set_network_name',
      data: {
        name,
        network
      },
      debounce: {
        delay: 500,
        key: `network_name:${network}`
      }
    };
  }

  return action;
}
