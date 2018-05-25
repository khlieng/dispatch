import { createSelector } from 'reselect';
import has from 'lodash/has';
import {
  findBreakpoints,
  messageHeight,
  linkify,
  timestamp,
  isChannel
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

export default createReducer(
  {},
  {
    [actions.ADD_MESSAGE](state, { server, tab, message }) {
      init(state, server, tab);
      state[server][tab].push(message);
    },

    [actions.ADD_MESSAGES](state, { server, tab, messages, prepend }) {
      if (prepend) {
        init(state, server, tab);
        state[server][tab].unshift(...messages);
      } else {
        messages.forEach(message => {
          init(state, server, message.tab || tab);
          state[server][message.tab || tab].push(message);
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

let nextID = 0;

function initMessage(message, tab, state) {
  if (message.time) {
    message.time = timestamp(new Date(message.time * 1000));
  } else {
    message.time = timestamp();
  }

  if (!message.id) {
    message.id = nextID;
    nextID++;
  }

  if (tab.charAt(0) === '#') {
    message.channel = true;
  }

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
