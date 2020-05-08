import sortBy from 'lodash/sortBy';
import createReducer from 'utils/createReducer';
import { updateSelection } from './tab';
import * as actions from './actions';

export const getPrivateChats = state => state.privateChats;

function open(state, server, nick) {
  if (!state[server]) {
    state[server] = [];
  }
  if (!state[server].includes(nick)) {
    state[server].push(nick);
    state[server] = sortBy(state[server], v => v.toLowerCase());
  }
}

export default createReducer(
  {},
  {
    [actions.OPEN_PRIVATE_CHAT](state, action) {
      open(state, action.server, action.nick);
    },

    [actions.CLOSE_PRIVATE_CHAT](state, { server, nick }) {
      const i = state[server]?.findIndex(n => n === nick);
      if (i !== -1) {
        state[server].splice(i, 1);
      }
    },

    [actions.PRIVATE_CHATS](state, { privateChats }) {
      privateChats.forEach(({ server, name }) => {
        if (!state[server]) {
          state[server] = [];
        }

        state[server].push(name);
      });
    },

    [actions.socket.PM](state, action) {
      if (action.from.indexOf('.') === -1) {
        open(state, action.server, action.from);
      }
    },

    [actions.DISCONNECT](state, { server }) {
      delete state[server];
    }
  }
);

export function openPrivateChat(server, nick) {
  return (dispatch, getState) => {
    if (!getState().privateChats[server]?.includes(nick)) {
      dispatch({
        type: actions.OPEN_PRIVATE_CHAT,
        server,
        nick,
        socket: {
          type: 'open_dm',
          data: { server, name: nick }
        }
      });
    }
  };
}

export function closePrivateChat(server, nick) {
  return dispatch => {
    dispatch({
      type: actions.CLOSE_PRIVATE_CHAT,
      server,
      nick,
      socket: {
        type: 'close_dm',
        data: { server, name: nick }
      }
    });
    dispatch(updateSelection());
  };
}
