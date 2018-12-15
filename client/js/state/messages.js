import { createSelector } from 'reselect';
import has from 'lodash/has';
import {
  findBreakpoints,
  messageHeight,
  linkify,
  timestamp,
  isChannel,
  formatDate
} from 'utils';
import createReducer from 'utils/createReducer';
import { getApp } from './app';
import { getSelectedTab } from './tab';
import * as actions from './actions';

export const getMessages = state => state.messages;

export const getSelectedMessages = createSelector(
  getSelectedTab,
  getMessages,
  (tab, messages) => {
    const target = tab.name || tab.server;
    if (has(messages, [tab.server, target])) {
      return messages[tab.server][target];
    }
    return [];
  }
);

export const getHasMoreMessages = createSelector(
  getSelectedMessages,
  messages => {
    const first = messages[0];
    return first && first.next;
  }
);

function init(state, server, tab) {
  if (!state[server]) {
    state[server] = {};
  }
  if (!state[server][tab]) {
    state[server][tab] = [];
  }
}

let nextID = 0;

function createDateMessage(date) {
  const message = {
    id: nextID,
    type: 'date',
    content: formatDate(date),
    height: 40
  };

  nextID++;

  return message;
}

function isSameDay(d1, d2) {
  return (
    d1.getDate() === d2.getDate() &&
    d1.getMonth() === d2.getMonth() &&
    d1.getFullYear() === d2.getFullYear()
  );
}

function reducerPrependMessages(messages, server, tab, state) {
  const msgs = [];

  for (let i = 0; i < messages.length; i++) {
    if (i > 0 && !isSameDay(messages[i - 1].date, messages[i].date)) {
      msgs.push(createDateMessage(messages[i].date));
    }
    msgs.push(messages[i]);
  }

  const m = state[server][tab];

  if (m.length > 0) {
    const lastNewMessage = msgs[msgs.length - 1];
    const firstMessage = m[0];
    if (
      firstMessage.date &&
      !isSameDay(firstMessage.date, lastNewMessage.date)
    ) {
      msgs.push(createDateMessage(firstMessage.date));
    }
  }

  m.unshift(...msgs);
}

function reducerAddMessage(message, server, tab, state) {
  const messages = state[server][tab];

  if (messages.length > 0) {
    const lastMessage = messages[messages.length - 1];
    if (lastMessage.date && !isSameDay(lastMessage.date, message.date)) {
      messages.push(createDateMessage(message.date));
    }
  }

  messages.push(message);
}

export default createReducer(
  {},
  {
    [actions.ADD_MESSAGE](state, { server, tab, message }) {
      init(state, server, tab);
      reducerAddMessage(message, server, tab, state);
    },

    [actions.ADD_MESSAGES](state, { server, tab, messages, prepend }) {
      if (prepend) {
        init(state, server, tab);
        reducerPrependMessages(messages, server, tab, state);
      } else {
        if (!messages[0].tab) {
          init(state, server, tab);
        }

        messages.forEach(message => {
          if (message.tab) {
            init(state, server, message.tab);
          }
          reducerAddMessage(message, server, message.tab || tab, state);
        });
      }
    },

    [actions.DISCONNECT](state, { server }) {
      delete state[server];
    },

    [actions.PART](state, { server, channels }) {
      channels.forEach(channel => delete state[server][channel]);
    },

    [actions.UPDATE_MESSAGE_HEIGHT](
      state,
      { wrapWidth, charWidth, windowWidth }
    ) {
      Object.keys(state).forEach(server =>
        Object.keys(state[server]).forEach(target =>
          state[server][target].forEach(message => {
            if (message.type === 'date') {
              return;
            }

            message.height = messageHeight(
              message,
              wrapWidth,
              charWidth,
              6 * charWidth,
              windowWidth
            );
          })
        )
      );
    },

    [actions.socket.SERVERS](state, { data }) {
      if (data) {
        data.forEach(({ host }) => {
          state[host] = {};
        });
      }
    }
  }
);

function initMessage(message, tab, state) {
  if (message.time) {
    message.date = new Date(message.time * 1000);
  } else {
    message.date = new Date();
  }

  message.time = timestamp(message.date);

  if (!message.id) {
    message.id = nextID;
    nextID++;
  }

  if (tab.charAt(0) === '#') {
    message.channel = true;
  }

  message.content = message.content || '';

  // Collapse multiple adjacent spaces into a single one
  message.content = message.content.replace(/\s\s+/g, ' ');

  if (message.content.indexOf('\x01ACTION') === 0) {
    const { from } = message;
    message.from = null;
    message.type = 'action';
    message.content = from + message.content.slice(7, -1);
  }

  const { wrapWidth, charWidth, windowWidth } = getApp(state);

  message.length = message.content.length;
  message.breakpoints = findBreakpoints(message.content);
  message.height = messageHeight(
    message,
    wrapWidth,
    charWidth,
    6 * charWidth,
    windowWidth
  );
  message.content = linkify(message.content);

  return message;
}

export function getMessageTab(server, to) {
  if (!to || to === '*' || (!isChannel(to) && to.indexOf('.') !== -1)) {
    return server;
  }
  return to;
}

export function fetchMessages() {
  return (dispatch, getState) => {
    const state = getState();
    const first = getSelectedMessages(state)[0];

    if (!first) {
      return;
    }

    const tab = state.tab.selected;
    if (isChannel(tab)) {
      dispatch({
        type: actions.FETCH_MESSAGES,
        socket: {
          type: 'fetch_messages',
          data: {
            server: tab.server,
            channel: tab.name,
            next: first.id
          }
        }
      });
    }
  };
}

export function addFetchedMessages(server, tab) {
  return {
    type: actions.ADD_FETCHED_MESSAGES,
    server,
    tab
  };
}

export function updateMessageHeight(wrapWidth, charWidth, windowWidth) {
  return {
    type: actions.UPDATE_MESSAGE_HEIGHT,
    wrapWidth,
    charWidth,
    windowWidth
  };
}

export function sendMessage(content, to, server) {
  return (dispatch, getState) => {
    const state = getState();

    dispatch({
      type: actions.ADD_MESSAGE,
      server,
      tab: to,
      message: initMessage(
        {
          from: state.servers[server].nick,
          content
        },
        to,
        state
      ),
      socket: {
        type: 'message',
        data: { content, to, server }
      }
    });
  };
}

export function addMessage(message, server, to) {
  const tab = getMessageTab(server, to);

  return (dispatch, getState) =>
    dispatch({
      type: actions.ADD_MESSAGE,
      server,
      tab,
      message: initMessage(message, tab, getState())
    });
}

export function addMessages(messages, server, to, prepend, next) {
  const tab = getMessageTab(server, to);

  return (dispatch, getState) => {
    const state = getState();

    if (next) {
      messages[0].id = next;
      messages[0].next = true;
    }

    messages.forEach(message =>
      initMessage(message, message.tab || tab, state)
    );

    dispatch({
      type: actions.ADD_MESSAGES,
      server,
      tab,
      messages,
      prepend
    });
  };
}

export function broadcast(message, server, channels) {
  return addMessages(
    channels.map(channel => ({
      tab: channel,
      content: message,
      type: 'info'
    })),
    server
  );
}

export function print(message, server, channel, type) {
  if (Array.isArray(message)) {
    return addMessages(
      message.map(line => ({
        content: line,
        type
      })),
      server,
      channel
    );
  }

  return addMessage(
    {
      content: message,
      type
    },
    server,
    channel
  );
}

export function inform(message, server, channel) {
  return print(message, server, channel, 'info');
}

export function runCommand(command, channel, server) {
  return {
    type: actions.COMMAND,
    command,
    channel,
    server
  };
}

export function raw(message, server) {
  return {
    type: actions.RAW,
    message,
    server,
    socket: {
      type: 'raw',
      data: { message, server }
    }
  };
}
