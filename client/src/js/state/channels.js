import { Map, List, Record } from 'immutable';
import { createSelector } from 'reselect';
import createReducer from '../util/createReducer';
import { getSelectedTab, updateSelection } from './tab';
import * as actions from './actions';

const User = Record({
  nick: null,
  renderName: null,
  mode: ''
});

const modePrefixes = [
  { mode: 'q', prefix: '~' }, // Owner
  { mode: 'a', prefix: '&' }, // Admin
  { mode: 'o', prefix: '@' }, // Op
  { mode: 'h', prefix: '%' }, // Halfop
  { mode: 'v', prefix: '+' }  // Voice
];

function updateRenderName(user) {
  for (let i = 0; i < modePrefixes.length; i++) {
    if (user.mode.indexOf(modePrefixes[i].mode) !== -1) {
      return user.set('renderName', `${modePrefixes[i].prefix}${user.nick}`);
    }
  }

  return user.set('renderName', user.nick);
}

function createUser(nick, mode) {
  return updateRenderName(new User({
    nick,
    renderName: nick,
    mode: mode || ''
  }));
}

function loadUser(nick) {
  let mode;

  for (let i = 0; i < modePrefixes.length; i++) {
    if (nick[0] === modePrefixes[i].prefix) {
      mode = modePrefixes[i].mode;
    }
  }

  if (mode) {
    return createUser(nick.slice(1), mode);
  }

  return createUser(nick);
}

function compareUsers(a, b) {
  a = a.renderName.toLowerCase();
  b = b.renderName.toLowerCase();

  for (let i = 0; i < modePrefixes.length; i++) {
    const prefix = modePrefixes[i].prefix;

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

export const getSelectedChannel = createSelector(
  getSelectedTab,
  getChannels,
  (tab, channels) => channels.getIn([tab.server, tab.name], Map())
);

export const getSelectedChannelUsers = createSelector(
  getSelectedChannel,
  channel => channel.get('users', List())
);

export default createReducer(Map(), {
  [actions.PART](state, action) {
    const { channels, server } = action;
    return state.withMutations(s => {
      channels.forEach(channel => s.deleteIn([server, channel]));
    });
  },

  [actions.socket.JOIN](state, action) {
    const { server, channels, user } = action;
    return state.updateIn([server, channels[0], 'users'], List(), users =>
      users.push(createUser(user)).sort(compareUsers)
    );
  },

  [actions.socket.PART](state, action) {
    const { server, channel, user } = action;
    if (state.hasIn([server, channel])) {
      return state.updateIn([server, channel, 'users'], users =>
        users.filter(u => u.nick !== user)
      );
    }
    return state;
  },

  [actions.socket.QUIT](state, action) {
    const { server, user } = action;
    return state.withMutations(s => {
      s.get(server).forEach((v, channel) => {
        s.updateIn([server, channel, 'users'], users => users.filter(u => u.nick !== user));
      });
    });
  },

  [actions.socket.NICK](state, { server, oldNick, newNick }) {
    return state.withMutations(s => {
      s.get(server).forEach((v, channel) => {
        s.updateIn([server, channel, 'users'], users => {
          const i = users.findIndex(user => user.nick === oldNick);
          if (i < 0) {
            return users;
          }

          return users.update(i,
            user => updateRenderName(user.set('nick', newNick))
          ).sort(compareUsers);
        });
      });
    });
  },

  [actions.socket.USERS](state, action) {
    const { server, channel, users } = action;
    return state.setIn([server, channel, 'users'],
      List(users.map(user => loadUser(user)).sort(compareUsers)));
  },

  [actions.socket.TOPIC](state, action) {
    const { server, channel, topic } = action;
    return state.setIn([server, channel, 'topic'], topic);
  },

  [actions.socket.MODE](state, action) {
    const { server, channel, user, remove, add } = action;

    return state.updateIn([server, channel, 'users'], users => {
      const i = users.findIndex(u => u.nick === user);
      if (i < 0) {
        return users;
      }

      return users.update(i, u => {
        let mode = u.mode;
        let j = remove.length;
        while (j--) {
          mode = mode.replace(remove[j], '');
        }

        return updateRenderName(u.set('mode', mode + add));
      }).sort(compareUsers);
    });
  },

  [actions.socket.CHANNELS](state, action) {
    if (!action.data) {
      return state;
    }

    return state.withMutations(s => {
      action.data.forEach(channel => {
        s.setIn([channel.server, channel.name], Map({
          users: List(),
          topic: channel.topic
        }));
      });
    });
  },

  [actions.socket.SERVERS](state, action) {
    if (!action.data) {
      return state;
    }

    return state.withMutations(s => {
      action.data.forEach(server => {
        if (!state.has(server.host)) {
          s.set(server.host, Map());
        }
      });
    });
  },

  [actions.CONNECT](state, action) {
    const { host } = action;

    if (!state.has(host)) {
      return state.set(host, Map());
    }

    return state;
  },

  [actions.DISCONNECT](state, action) {
    return state.delete(action.server);
  }
});

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
  return dispatch => {
    dispatch({
      type: actions.PART,
      channels,
      server,
      socket: {
        type: 'part',
        data: { channels, server }
      }
    });
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
