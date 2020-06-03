import React from 'react';
import { createSelector } from 'reselect';
import has from 'lodash/has';
import {
  findBreakpoints,
  messageHeight,
  linkify,
  timestamp,
  isChannel,
  formatDate,
  unix
} from 'utils';
import stringToRGB from 'utils/color';
import colorify from 'utils/colorify';
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

const collapsedEvents = ['join', 'part', 'quit'];

function shouldCollapse(msg1, msg2) {
  return (
    msg1.events &&
    msg2.events &&
    collapsedEvents.indexOf(msg1.events[0].type) !== -1 &&
    collapsedEvents.indexOf(msg2.events[0].type) !== -1
  );
}

const eventVerbs = {
  join: 'joined the channel',
  part: 'left the channel',
  quit: 'quit'
};

function renderNick(nick, type = '') {
  const style = {
    color: stringToRGB(nick),
    fontWeight: 400
  };

  return (
    <span className="message-sender" style={style} key={`${nick} ${type}`}>
      {nick}
    </span>
  );
}

function renderMore(count, type) {
  return (
    <span
      className="message-events-more"
      key={`more ${type}`}
    >{`${count} more`}</span>
  );
}

function renderEvent(event, type, nicks) {
  const ending = eventVerbs[type];

  if (nicks.length === 1) {
    event.push(renderNick(nicks[0], type));
    event.push(` ${ending}`);
  }
  if (nicks.length === 2) {
    event.push(renderNick(nicks[0], type));
    event.push(' and ');
    event.push(renderNick(nicks[1], type));
    event.push(` ${ending}`);
  }
  if (nicks.length > 2) {
    event.push(renderNick(nicks[0], type));
    event.push(', ');
    event.push(renderNick(nicks[1], type));
    event.push(' and ');
    event.push(renderMore(nicks.length - 2, type));
    event.push(` ${ending}`);
  }
}

function renderEvents(events) {
  const first = events[0];
  if (first.type === 'nick') {
    const [oldNick, newNick] = first.params;

    return [renderNick(oldNick), ' changed nick to ', renderNick(newNick)];
  }
  if (first.type === 'topic') {
    const [nick, newTopic] = first.params;
    const topic = colorify(linkify(newTopic));

    if (!topic) {
      return [renderNick(nick), ' cleared the topic'];
    }

    const result = [renderNick(nick), ' changed the topic to: '];

    if (Array.isArray(topic)) {
      result.push(...topic);
    } else {
      result.push(topic);
    }

    return result;
  }

  const byType = {};
  for (let i = events.length - 1; i >= 0; i--) {
    const event = events[i];
    const [nick] = event.params;

    if (!byType[event.type]) {
      byType[event.type] = [nick];
    } else if (byType[event.type].indexOf(nick) === -1) {
      byType[event.type].push(nick);
    }
  }

  const result = [];

  if (byType.join) {
    renderEvent(result, 'join', byType.join);
  }

  if (byType.part) {
    if (result.length > 1) {
      result[result.length - 1] += ', ';
    }
    renderEvent(result, 'part', byType.part);
  }

  if (byType.quit) {
    if (result.length > 1) {
      result[result.length - 1] += ', ';
    }
    renderEvent(result, 'quit', byType.quit);
  }

  return result;
}

let nextID = 0;

function initMessage(
  state,
  message,
  server,
  tab,
  wrapWidth,
  charWidth,
  windowWidth,
  prepend
) {
  const messages = state[server][tab];

  if (messages.length > 0 && !prepend) {
    const lastMessage = messages[messages.length - 1];
    if (shouldCollapse(lastMessage, message)) {
      lastMessage.events.push(message.events[0]);
      lastMessage.content = renderEvents(lastMessage.events);

      [lastMessage.breakpoints, lastMessage.length] = findBreakpoints(
        lastMessage.content
      );
      lastMessage.height = messageHeight(
        lastMessage,
        wrapWidth,
        charWidth,
        6 * charWidth,
        windowWidth
      );

      return false;
    }
  }

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

  if (message.events) {
    message.type = 'info';
    message.content = renderEvents(message.events);
  } else {
    message.content = message.content || '';
    // Collapse multiple adjacent spaces into a single one
    message.content = message.content.replace(/\s\s+/g, ' ');

    if (message.content.indexOf('\x01ACTION') === 0) {
      const { from } = message;
      message.from = null;
      message.type = 'action';
      message.content = from + message.content.slice(7, -1);
    }
  }

  if (!message.events) {
    message.content = colorify(linkify(message.content));
  }

  [message.breakpoints, message.length] = findBreakpoints(message.content);
  message.height = messageHeight(
    message,
    wrapWidth,
    charWidth,
    6 * charWidth,
    windowWidth
  );
  message.indent = 6 * charWidth;

  return true;
}

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

function reducerPrependMessages(
  state,
  messages,
  server,
  tab,
  wrapWidth,
  charWidth,
  windowWidth
) {
  const msgs = [];

  for (let i = 0; i < messages.length; i++) {
    const message = messages[i];
    initMessage(
      state,
      message,
      server,
      tab,
      wrapWidth,
      charWidth,
      windowWidth,
      true
    );

    if (i > 0 && !isSameDay(messages[i - 1].date, message.date)) {
      msgs.push(createDateMessage(message.date));
    }
    msgs.push(message);
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
    [actions.ADD_MESSAGE](
      state,
      { server, tab, message, wrapWidth, charWidth, windowWidth }
    ) {
      init(state, server, tab);

      const shouldAdd = initMessage(
        state,
        message,
        server,
        tab,
        wrapWidth,
        charWidth,
        windowWidth
      );
      if (shouldAdd) {
        reducerAddMessage(message, server, tab, state);
      }
    },

    [actions.ADD_MESSAGES](
      state,
      { server, tab, messages, prepend, wrapWidth, charWidth, windowWidth }
    ) {
      if (prepend) {
        init(state, server, tab);
        reducerPrependMessages(
          state,
          messages,
          server,
          tab,
          wrapWidth,
          charWidth,
          windowWidth
        );
      } else {
        if (!messages[0].tab) {
          init(state, server, tab);
        }

        messages.forEach(message => {
          if (message.tab) {
            init(state, server, message.tab);
          }

          const shouldAdd = initMessage(
            state,
            message,
            server,
            message.tab || tab,
            wrapWidth,
            charWidth,
            windowWidth
          );
          if (shouldAdd) {
            reducerAddMessage(message, server, message.tab || tab, state);
          }
        });
      }
    },

    [actions.DISCONNECT](state, { server }) {
      delete state[server];
    },

    [actions.PART](state, { server, channels }) {
      channels.forEach(channel => delete state[server][channel]);
    },

    [actions.CLOSE_PRIVATE_CHAT](state, { server, nick }) {
      delete state[server][nick];
    },

    [actions.socket.CHANNEL_FORWARD](state, { server, old }) {
      if (state[server]) {
        delete state[server][old];
      }
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
    if (tab.name) {
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
    const { wrapWidth, charWidth, windowWidth } = getApp(state);

    dispatch({
      type: actions.ADD_MESSAGE,
      server,
      tab: to,
      message: {
        from: state.servers[server].nick,
        content
      },
      wrapWidth,
      charWidth,
      windowWidth,
      socket: {
        type: 'message',
        data: { content, to, server }
      }
    });
  };
}

export function addMessage(message, server, to) {
  const tab = getMessageTab(server, to);

  return (dispatch, getState) => {
    const { wrapWidth, charWidth, windowWidth } = getApp(getState());

    dispatch({
      type: actions.ADD_MESSAGE,
      server,
      tab,
      message,
      wrapWidth,
      charWidth,
      windowWidth
    });
  };
}

export function addMessages(messages, server, to, prepend, next) {
  const tab = getMessageTab(server, to);

  return (dispatch, getState) => {
    const state = getState();

    if (next) {
      messages[0].id = next;
      messages[0].next = true;
    }

    const { wrapWidth, charWidth, windowWidth } = getApp(state);

    dispatch({
      type: actions.ADD_MESSAGES,
      server,
      tab,
      messages,
      prepend,
      wrapWidth,
      charWidth,
      windowWidth
    });
  };
}

export function addEvent(server, tab, type, ...params) {
  return addMessage(
    {
      type: 'info',
      events: [
        {
          type,
          params,
          time: unix()
        }
      ]
    },
    server,
    tab
  );
}

export function broadcastEvent(server, channels, type, ...params) {
  const now = unix();

  return addMessages(
    channels.map(channel => ({
      type: 'info',
      tab: channel,
      events: [
        {
          type,
          params,
          time: now
        }
      ]
    })),
    server
  );
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
