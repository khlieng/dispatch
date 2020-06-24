import { createSelector } from 'reselect';
import get from 'lodash/get';
import sortBy from 'lodash/sortBy';
import createReducer from 'utils/createReducer';
import { trimPrefixChar, find, findIndex } from 'utils';
import { getSelectedTab, updateSelection } from './tab';
import * as actions from './actions';

const modePrefixes = [
  { mode: 'q', prefix: '~' }, // Owner
  { mode: 'a', prefix: '&' }, // Admin
  { mode: 'o', prefix: '@' }, // Op
  { mode: 'h', prefix: '%' }, // Halfop
  { mode: 'v', prefix: '+' } // Voice
];

function getRenderName(user) {
  for (let i = 0; i < modePrefixes.length; i++) {
    if (user.mode.indexOf(modePrefixes[i].mode) !== -1) {
      return `${modePrefixes[i].prefix}${user.nick}`;
    }
  }

  return user.nick;
}

function createUser(nick, mode) {
  const user = {
    nick,
    mode: mode || ''
  };
  user.renderName = getRenderName(user);
  return user;
}

function loadUser(nick) {
  let mode;

  for (let i = 0; i < modePrefixes.length; i++) {
    if (nick[0] === modePrefixes[i].prefix) {
      ({ mode } = modePrefixes[i]);
    }
  }

  if (mode) {
    return createUser(nick.slice(1), mode);
  }

  return createUser(nick);
}

function removeUser(users, nick) {
  const i = findIndex(users, u => u.nick === nick);
  if (i !== -1) {
    users.splice(i, 1);
  }
}

function init(state, network, channel) {
  if (!state[network]) {
    state[network] = {};
  }
  if (channel && !state[network][channel]) {
    state[network][channel] = {
      name: channel,
      users: [],
      joined: false
    };
  }
  return state[network][channel];
}

export function compareUsers(a, b) {
  a = a.renderName.toLowerCase();
  b = b.renderName.toLowerCase();

  for (let i = 0; i < modePrefixes.length; i++) {
    const { prefix } = modePrefixes[i];

    if (a[0] === prefix && b[0] !== prefix) {
      return -1;
    }
    if (b[0] === prefix && a[0] !== prefix) {
      return 1;
    }
  }

  if (a < b) {
    return -1;
  }
  if (a > b) {
    return 1;
  }
  return 0;
}

export const getChannels = state => state.channels;

export const getSortedChannels = createSelector(getChannels, channels =>
  sortBy(
    Object.keys(channels).map(network => ({
      address: network,
      channels: sortBy(channels[network], channel =>
        trimPrefixChar(channel.name, '#').toLowerCase()
      )
    })),
    network => network.address.toLowerCase()
  )
);

export const getSelectedChannel = createSelector(
  getSelectedTab,
  getChannels,
  (tab, channels) => get(channels, [tab.network, tab.name])
);

export const getSelectedChannelUsers = createSelector(
  getSelectedChannel,
  channel => {
    if (channel) {
      return channel.users.concat().sort(compareUsers);
    }
    return [];
  }
);

export default createReducer(
  {},
  {
    [actions.JOIN](state, { network, channels }) {
      channels.forEach(channel => init(state, network, channel));
    },

    [actions.PART](state, { network, channels }) {
      channels.forEach(channel => delete state[network][channel]);
    },

    [actions.socket.JOIN](state, { network, channels, user }) {
      const channel = channels[0];
      const chan = init(state, network, channel);
      chan.name = channel;
      chan.joined = true;
      chan.users.push(createUser(user));
    },

    [actions.socket.CHANNEL_FORWARD](state, action) {
      init(state, action.network, action.new);
      delete state[action.network][action.old];
    },

    [actions.socket.PART](state, { network, channel, user }) {
      if (state[network][channel]) {
        removeUser(state[network][channel].users, user);
      }
    },

    [actions.socket.QUIT](state, { network, user }) {
      Object.keys(state[network]).forEach(channel => {
        removeUser(state[network][channel].users, user);
      });
    },

    [actions.KICKED](state, { network, channel, user, self }) {
      const chan = state[network][channel];
      if (self) {
        chan.joined = false;
        chan.users = [];
      } else {
        removeUser(chan.users, user);
      }
    },

    [actions.socket.NICK](state, { network, oldNick, newNick }) {
      Object.keys(state[network]).forEach(channel => {
        const user = find(
          state[network][channel].users,
          u => u.nick === oldNick
        );
        if (user) {
          user.nick = newNick;
          user.renderName = getRenderName(user);
        }
      });
    },

    [actions.socket.MODE](state, { network, channel, user, remove, add }) {
      const u = find(state[network][channel].users, v => v.nick === user);
      if (u) {
        if (remove) {
          let j = remove.length;
          while (j--) {
            u.mode = u.mode.replace(remove[j], '');
          }
        }

        if (add) {
          u.mode += add;
        }

        u.renderName = getRenderName(u);
      }
    },

    [actions.socket.TOPIC](state, { network, channel, topic }) {
      state[network][channel].topic = topic;
    },

    [actions.socket.USERS](state, { network, channel, users }) {
      state[network][channel].users = users.map(nick => loadUser(nick));
    },

    [actions.INIT](state, { networks, channels, users }) {
      if (networks) {
        networks.forEach(({ host }) => init(state, host));
      }

      if (channels) {
        channels.forEach(({ network, name, topic, joined }) => {
          const chan = init(state, network, name);
          chan.joined = joined;
          chan.topic = topic;
        });
      }

      if (users) {
        state[users.network][users.channel].users = users.users.map(nick =>
          loadUser(nick)
        );
      }
    },

    [actions.CONNECT](state, { host }) {
      init(state, host);
    },

    [actions.DISCONNECT](state, { network }) {
      delete state[network];
    }
  }
);

export function join(channels, network, selectFirst = true) {
  return {
    type: actions.JOIN,
    channels,
    network,
    selectFirst,
    socket: {
      type: 'join',
      data: { channels, network }
    }
  };
}

export function part(channels, network) {
  return (dispatch, getState) => {
    const action = {
      type: actions.PART,
      channels,
      network
    };

    const state = getState().channels[network];
    const joined = channels.filter(c => state[c] && state[c].joined);

    if (joined.length > 0) {
      action.socket = {
        type: 'part',
        data: {
          channels: joined,
          network
        }
      };
    }

    dispatch(action);
    dispatch(updateSelection());
  };
}

export function invite(user, channel, network) {
  return {
    type: actions.INVITE,
    user,
    channel,
    network,
    socket: {
      type: 'invite',
      data: { user, channel, network }
    }
  };
}

export function kick(user, channel, network) {
  return {
    type: actions.KICK,
    user,
    channel,
    network,
    socket: {
      type: 'kick',
      data: { user, channel, network }
    }
  };
}

export function kicked(network, channel, user) {
  return (dispatch, getState) => {
    const nick = getState().networks[network]?.nick;

    dispatch({
      type: actions.KICKED,
      network,
      channel,
      user,
      self: nick === user
    });
  };
}

export function setTopic(topic, channel, network) {
  return {
    type: actions.SET_TOPIC,
    topic,
    channel,
    network,
    socket: {
      type: 'topic',
      data: { topic, channel, network }
    }
  };
}
