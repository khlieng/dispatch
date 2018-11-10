import sortBy from 'lodash/sortBy';
import { findIndex } from 'utils';
import createReducer from 'utils/createReducer';
import { updateSelection } from './tab';
import * as actions from './actions';

export const getPrivateChats = state => state.privateChats;

function open(state, server, nick) {
  if (!state[server]) {
    state[server] = [];
  }
  if (findIndex(state[server], n => n === nick) === -1) {
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
      const i = findIndex(state[server], n => n === nick);
      if (i !== -1) {
        state[server].splice(i, 1);
      }
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
  return {
    type: actions.OPEN_PRIVATE_CHAT,
    server,
    nick
  };
}

export function closePrivateChat(server, nick) {
  return dispatch => {
    dispatch({
      type: actions.CLOSE_PRIVATE_CHAT,
      server,
      nick
    });
    dispatch(updateSelection());
  };
}
