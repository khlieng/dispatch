import { createSelector } from 'reselect';
import get from 'lodash/get';
import sortBy from 'lodash/sortBy';
import createReducer from 'utils/createReducer';
import { find, findIndex } from 'utils';
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

function init(state, server, channel) {
  if (!state[server]) {
    state[server] = {};
  }
  if (channel && !state[server][channel]) {
    state[server][channel] = { name: channel, users: [], joined: false };
  }
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

export const getSortedChannels = createSelector(
  getChannels,
  channels =>
    sortBy(
      Object.keys(channels).map(server => ({
        address: server,
        channels: sortBy(Object.keys(channels[server]), channel =>
          channel.toLowerCase()
        )
      })),
      server => server.address.toLowerCase()
    )
);

export const getSelectedChannel = createSelector(
  getSelectedTab,
  getChannels,
  (tab, channels) => get(channels, [tab.server, tab.name])
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
    [actions.JOIN](state, { server, channels }) {
      channels.forEach(channel => init(state, server, channel));
    },

    [actions.PART](state, { server, channels }) {
      channels.forEach(channel => delete state[server][channel]);
    },

    [actions.socket.JOIN](state, { server, channels, user }) {
      const channel = channels[0];
      init(state, server, channel);
      state[server][channel].name = channel;
      state[server][channel].joined = true;
      state[server][channel].users.push(createUser(user));
    },

    [actions.socket.CHANNEL_FORWARD](state, action) {
      init(state, action.server, action.new);
      delete state[action.server][action.old];
    },

    [actions.socket.PART](state, { server, channel, user }) {
      if (state[server][channel]) {
        removeUser(state[server][channel].users, user);
      }
    },

    [actions.socket.QUIT](state, { server, user }) {
      Object.keys(state[server]).forEach(channel => {
        removeUser(state[server][channel].users, user);
      });
    },

    [actions.socket.NICK](state, { server, oldNick, newNick }) {
      Object.keys(state[server]).forEach(channel => {
        const user = find(
          state[server][channel].users,
          u => u.nick === oldNick
        );
        if (user) {
          user.nick = newNick;
          user.renderName = getRenderName(user);
        }
      });
    },

    [actions.socket.USERS](state, { server, channel, users }) {
      state[server][channel].users = users.map(nick => loadUser(nick));
    },

    [actions.socket.TOPIC](state, { server, channel, topic }) {
      state[server][channel].topic = topic;
    },

    [actions.socket.MODE](state, { server, channel, user, remove, add }) {
      const u = find(state[server][channel].users, v => v.nick === user);
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

    [actions.socket.CHANNELS](state, { data }) {
      if (data) {
        data.forEach(({ server, name, topic }) => {
          init(state, server, name);
          state[server][name].joined = true;
          state[server][name].topic = topic;
        });
      }
    },

    [actions.socket.SERVERS](state, { data }) {
      if (data) {
        data.forEach(({ host }) => init(state, host));
      }
    },

    [actions.CONNECT](state, { host }) {
      init(state, host);
    },

    [actions.DISCONNECT](state, { server }) {
      delete state[server];
    }
  }
);

export function join(channels, server) {
  return {
    type: actions.JOIN,
    channels,
    server,
    socket: {
      type: 'join',
      data: { channels, server }
    }
  };
}

export function part(channels, server) {
  return (dispatch, getState) => {
    const action = {
      type: actions.PART,
      channels,
      server
    };

    const state = getState().channels[server];
    const joined = channels.filter(c => state[c] && state[c].joined);

    if (joined.length > 0) {
      action.socket = {
        type: 'part',
        data: {
          channels: joined,
          server
        }
      };
    }

    dispatch(action);
    dispatch(updateSelection());
  };
}

export function invite(user, channel, server) {
  return {
    type: actions.INVITE,
    user,
    channel,
    server,
    socket: {
      type: 'invite',
      data: { user, channel, server }
    }
  };
}

export function kick(user, channel, server) {
  return {
    type: actions.KICK,
    user,
    channel,
    server,
    socket: {
      type: 'kick',
      data: { user, channel, server }
    }
  };
}

export function setTopic(topic, channel, server) {
  return {
    type: actions.SET_TOPIC,
    topic,
    channel,
    server,
    socket: {
      type: 'topic',
      data: { topic, channel, server }
    }
  };
}
