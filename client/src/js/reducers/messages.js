import { List, Map, Record } from 'immutable';
import { createSelector } from 'reselect';
import createReducer from '../util/createReducer';
import { messageHeight } from '../util';
import * as actions from '../actions';
import { getSelectedTab } from './tab';

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

export const getMessages = state => state.messages;

export const getSelectedMessages = createSelector(
  getSelectedTab,
  getMessages,
  (tab, messages) => messages.getIn([tab.server, tab.name || tab.server], List())
);

export default createReducer(Map(), {
  [actions.SEND_MESSAGE]: addMessage,
  [actions.ADD_MESSAGE]: addMessage,

  [actions.ADD_MESSAGES](state, { server, tab, messages, prepend }) {
    return state.withMutations(s => {
      if (prepend) {
        for (let i = messages.length - 1; i >= 0; i--) {
          s.updateIn([server, tab], List(), list => list.unshift(new Message(messages[i])));
        }
      } else {
        messages.forEach(message =>
          s.updateIn([server, tab], List(), list => list.push(new Message(message)))
        );
      }
    });
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
