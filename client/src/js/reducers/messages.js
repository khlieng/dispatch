import { List, Map, Record } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from '../actions';

const Message = Record({
  id: null,
  server: null,
  from: null,
  to: null,
  message: '',
  time: null,
  type: null,
  lines: []
});

function addMessage(state, message) {
  let dest = message.to || message.from;
  if (message.from && message.from.indexOf('.') !== -1) {
    dest = message.server;
  }

  if (message.message.indexOf('\x01ACTION') === 0) {
    const from = message.from;
    message.from = null;
    message.type = 'action';
    message.message = from + message.message.slice(7);
  }

  return state.updateIn([message.server, dest], List(), list => list.push(new Message(message)));
}

export default createReducer(Map(), {
  [actions.SEND_MESSAGE](state, action) {
    return addMessage(state, action);
  },

  [actions.ADD_MESSAGE](state, action) {
    return addMessage(state, action.message);
  },

  [actions.ADD_MESSAGES](state, action) {
    return state.withMutations(s =>
      action.messages.forEach(message =>
        addMessage(s, message)
      )
    );
  },
/*
  [actions.SOCKET_MESSAGE](state, action) {
    return addMessage(state, action);
  },

  [actions.SOCKET_PM](state, action) {
    return addMessage(state, action);
  },
*/
  [actions.DISCONNECT](state, action) {
    return state.delete(action.server);
  },

  [actions.PART](state, action) {
    return state.withMutations(s =>
      action.channels.forEach(channel =>
        s.deleteIn([action.server, channel])
      )
    );
  }
});
