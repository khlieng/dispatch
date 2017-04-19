import { List, Map, Record } from 'immutable';
import createReducer from '../util/createReducer';
import { messageHeight } from '../util';
import * as actions from '../actions';

const Message = Record({
  id: null,
  from: null,
  content: '',
  time: null,
  type: null,
  channel: false,
  height: 0,
  length: 0,
  breakpoints: null
});

function addMessage(state, { server, tab, message }) {
  return state.updateIn([server, tab], List(), list => list.push(new Message(message)));
}

export default createReducer(Map(), {
  [actions.SEND_MESSAGE]: addMessage,
  [actions.ADD_MESSAGE]: addMessage,

  [actions.ADD_MESSAGES](state, { server, tab, messages }) {
    return state.withMutations(s =>
      messages.forEach(message =>
        s.updateIn([server, tab], List(), list => list.push(new Message(message)))
      )
    );
  },

  [actions.DISCONNECT](state, { server }) {
    return state.delete(server);
  },

  [actions.PART](state, { server, channels }) {
    return state.withMutations(s =>
      channels.forEach(channel =>
        s.deleteIn([server, channel])
      )
    );
  },

  [actions.UPDATE_MESSAGE_HEIGHT](state, { wrapWidth, charWidth }) {
    return state.withMutations(s =>
      s.forEach((server, serverKey) =>
        server.forEach((target, targetKey) =>
          target.forEach((message, index) => s.setIn([serverKey, targetKey, index, 'height'],
            messageHeight(message, wrapWidth, charWidth, 6 * charWidth))
          )
        )
      )
    );
  }
});
