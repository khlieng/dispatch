import { Map, List, Record } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from '../actions';

const User = Record({
  nick: null,
  renderName: null,
  mode: ''
});

function updateRenderName(user) {
  let name = user.nick;

  if (user.mode.indexOf('o') !== -1) {
    name = '@' + name;
  } else if (user.mode.indexOf('v') !== -1) {
    name = '+' + name;
  }

  return user.set('renderName', name);
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

  if (nick[0] === '@') {
    mode = 'o';
  } else if (nick[0] === '+') {
    mode = 'v';
  }

  if (mode) {
    return createUser(nick.slice(1), mode);
  }

  return createUser(nick, mode);
}

function compareUsers(a, b) {
  a = a.renderName.toLowerCase();
  b = b.renderName.toLowerCase();

  if (a[0] === '@' && b[0] !== '@') {
    return -1;
  }
  if (b[0] === '@' && a[0] !== '@') {
    return 1;
  }
  if (a[0] === '+' && b[0] !== '+') {
    return -1;
  }
  if (b[0] === '+' && a[0] !== '+') {
    return 1;
  }
  if (a < b) {
    return -1;
  }
  if (a > b) {
    return 1;
  }
  return 0;
}

export default createReducer(Map(), {
  [actions.PART](state, action) {
    const { channels, server } = action;
    return state.withMutations(s => {
      channels.forEach(channel => s.deleteIn([server, channel]));
    });
  },

  [actions.SOCKET_JOIN](state, action) {
    const { server, channels, user } = action;
    return state.updateIn([server, channels[0], 'users'], List(), users => {
      return users.push(createUser(user)).sort(compareUsers);
    });
  },

  [actions.SOCKET_PART](state, action) {
    const { server, channels, user } = action;
    const channel = channels[0];
    if (state.hasIn([server, channel])) {
      return state.updateIn([server, channel, 'users'], users =>
        users.filter(u => u.nick !== user)
      );
    }
    return state;
  },

  [actions.SOCKET_QUIT](state, action) {
    const { server, user } = action;
    return state.withMutations(s => {
      s.get(server).forEach((v, channel) => {
        s.updateIn([server, channel, 'users'], users => users.filter(u => u.nick !== user));
      });
    });
  },

  [actions.SOCKET_NICK](state, action) {
    const { server, channels } = action;
    return state.withMutations(s => {
      channels.forEach(channel => {
        s.updateIn([server, channel, 'users'], users => {
          const i = users.findIndex(user => user.nick === action.old);
          return users.update(i, user => {
            return updateRenderName(user.set('nick', action.new));
          }).sort(compareUsers);
        });
      });
    });
  },

  [actions.SOCKET_USERS](state, action) {
    const { server, channel, users } = action;
    return state.setIn([server, channel, 'users'],
      List(users.map(user => loadUser(user)).sort(compareUsers)));
  },

  [actions.SOCKET_TOPIC](state, action) {
    const { server, channel, topic } = action;
    return state.setIn([server, channel, 'topic'], topic);
  },

  [actions.SOCKET_MODE](state, action) {
    const { server, channel, user, remove, add } = action;

    const i = state.getIn([server, channel, 'users']).findIndex(u => u.nick === user);
    return state
      .updateIn([server, channel, 'users', i], u => {
        let mode = u.mode;
        let j = remove.length;
        while (j--) {
          mode = mode.replace(remove[j], '');
        }

        return updateRenderName(u.set('mode', mode + add));
      })
      .updateIn([server, channel, 'users'], users => users.sort(compareUsers));
  },

  [actions.SOCKET_CHANNELS](state, action) {
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

  [actions.SOCKET_SERVERS](state, action) {
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
