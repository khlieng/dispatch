import { Set, Map } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from '../actions';

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

  [actions.SOCKET_PM](state, action) {
    if (action.from.indexOf('.') === -1) {
      return open(state, action.server, action.from);
    }

    return state;
  },

  [actions.DISCONNECT](state, action) {
    return state.delete(action.server);
  }
});
