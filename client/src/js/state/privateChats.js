import { Set, Map } from 'immutable';
import createReducer from '../util/createReducer';
import { updateSelection } from './tab';
import * as actions from './actions';

export const getPrivateChats = state => state.privateChats;

function open(state, server, nick) {
  return state.update(server, Set(), chats => chats.add(nick));
}

export default createReducer(Map(), {
  [actions.OPEN_PRIVATE_CHAT](state, action) {
    return open(state, action.server, action.nick);
  },

  [actions.CLOSE_PRIVATE_CHAT](state, action) {
    return state.update(action.server, chats => chats.delete(action.nick));
  },

  [actions.socket.PM](state, action) {
    if (action.from.indexOf('.') === -1) {
      return open(state, action.server, action.from);
    }

    return state;
  },

  [actions.DISCONNECT](state, action) {
    return state.delete(action.server);
  }
});

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
