import { List, Map, Record } from 'immutable';
import createReducer from '../util/createReducer';
import { messageHeight } from '../util';
import * as actions from '../actions';

const Message = Record({
  id: null,
  server: null,
  from: null,
  to: null,
  message: '',
  time: null,
  type: null,
  height: 0,
  channel: false
});

function addMessage(state, message) {
  if (message.message.indexOf('\x01ACTION') === 0) {
    const from = message.from;
    message.from = null;
    message.type = 'action';
    message.message = from + message.message.slice(7, -1);
  }

  return state.updateIn([message.server, message.dest], List(),
    list => list.push(new Message(message)));
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

  [actions.DISCONNECT](state, action) {
    return state.delete(action.server);
  },

  [actions.PART](state, action) {
    return state.withMutations(s =>
      action.channels.forEach(channel =>
        s.deleteIn([action.server, channel])
      )
    );
  },

  [actions.UPDATE_MESSAGE_HEIGHT](state, action) {
    return state.withMutations(s =>
      s.forEach((server, serverKey) =>
        server.forEach((target, targetKey) =>
          target.forEach((message, index) => s.setIn([serverKey, targetKey, index, 'height'],
            messageHeight(message, action.wrapWidth, action.charWidth, 6 * action.charWidth))
          )
        )
      )
    );
  }
});
