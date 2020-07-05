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
    const target = tab.name || tab.network;
    if (has(messages, [tab.network, target])) {
      return messages[tab.network][target];
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

function init(state, network, tab) {
  if (!state[network]) {
    state[network] = {};
  }
  if (!state[network][tab]) {
    state[network][tab] = [];
  }
}

function initNetworks(state, networks = []) {
  networks.forEach(({ host }) => {
    state[host] = {};
  });
}

const collapsedEvents = ['join', 'part', 'quit', 'nick'];

function shouldCollapse(msg1, msg2) {
  return (
    msg1.events &&
    msg2.events &&
    collapsedEvents.indexOf(msg1.events[0].type) !== -1 &&
    collapsedEvents.indexOf(msg2.events[0].type) !== -1
  );
}

const blocks = {
  nick: nick => ({ type: 'nick', text: nick }),
  text: text => ({ type: 'text', text }),
  events: count => ({ type: 'events', text: `${count} more` })
};

const eventVerbs = {
  join: 'joined',
  part: 'left',
  quit: 'quit'
};

function renderEvent(result, type, events) {
  const ending = eventVerbs[type];

  if (result.length > 1) {
    result[result.length - 1].text += ', ';
  }

  if (events.length === 1) {
    result.push(blocks.nick(events[0][0]));
    result.push(blocks.text(` ${ending}`));
  } else if (events.length === 2) {
    result.push(blocks.nick(events[0][0]));
    result.push(blocks.text(' and '));
    result.push(blocks.nick(events[1][0]));
    result.push(blocks.text(` ${ending}`));
  } else if (events.length > 2) {
    result.push(blocks.nick(events[0][0]));
    result.push(blocks.text(', '));
    result.push(blocks.nick(events[1][0]));
    result.push(blocks.text(' and '));
    result.push(blocks.events(events.length - 2));
    result.push(blocks.text(` ${ending}`));
  }
}

function renderEvents(events) {
  const first = events[0];

  if (first.type === 'kick') {
    const [kicked, by] = first.params;

    return [blocks.nick(by), blocks.text(' kicked '), blocks.nick(kicked)];
  }

  if (first.type === 'topic') {
    const [nick, topic] = first.params;

    if (!topic) {
      return [blocks.nick(nick), blocks.text(' cleared the topic')];
    }

    return [
      blocks.nick(nick),
      blocks.text(' changed the topic to: '),
      ...colorify(linkify(topic))
    ];
  }

  const byType = {};
  for (let i = events.length - 1; i >= 0; i--) {
    const event = events[i];
    const [nick] = event.params;

    if (!byType[event.type]) {
      byType[event.type] = [event.params];
    } else if (byType[event.type].indexOf(nick) === -1) {
      byType[event.type].push(event.params);
    }
  }

  const result = [];

  if (byType.join) {
    renderEvent(result, 'join', byType.join);
  }

  if (byType.part) {
    renderEvent(result, 'part', byType.part);
  }

  if (byType.quit) {
    renderEvent(result, 'quit', byType.quit);
  }

  if (byType.nick) {
    if (result.length > 1) {
      result[result.length - 1].text += ', ';
    }

    const [oldNick, newNick] = byType.nick[0];

    result.push(blocks.nick(oldNick));
    result.push(blocks.text(' changed nick to '));
    result.push(blocks.nick(newNick));

    if (byType.nick.length > 1) {
      result.push(blocks.text(' and '));
      result.push(blocks.events(byType.nick.length - 1));
      result.push(blocks.text(' changed nick'));
    }
  }

  return result;
}

let nextID = 0;

function initMessage(
  state,
  message,
  network,
  tab,
  wrapWidth,
  charWidth,
  windowWidth,
  prepend
) {
  const messages = state[network][tab];

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
  network,
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
      network,
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

  const m = state[network][tab];

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

function reducerAddMessage(message, network, tab, state) {
  const messages = state[network][tab];

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
      { network, tab, message, wrapWidth, charWidth, windowWidth }
    ) {
      init(state, network, tab);

      const shouldAdd = initMessage(
        state,
        message,
        network,
        tab,
        wrapWidth,
        charWidth,
        windowWidth
      );
      if (shouldAdd) {
        reducerAddMessage(message, network, tab, state);
      }
    },

    [actions.ADD_MESSAGES](
      state,
      { network, tab, messages, prepend, wrapWidth, charWidth, windowWidth }
    ) {
      if (prepend) {
        init(state, network, tab);
        reducerPrependMessages(
          state,
          messages,
          network,
          tab,
          wrapWidth,
          charWidth,
          windowWidth
        );
      } else {
        if (!messages[0].tab) {
          init(state, network, tab);
        }

        messages.forEach(message => {
          if (message.tab) {
            init(state, network, message.tab);
          }

          const shouldAdd = initMessage(
            state,
            message,
            network,
            message.tab || tab,
            wrapWidth,
            charWidth,
            windowWidth
          );
          if (shouldAdd) {
            reducerAddMessage(message, network, message.tab || tab, state);
          }
        });
      }
    },

    [actions.DISCONNECT](state, { network }) {
      delete state[network];
    },

    [actions.PART](state, { network, channels }) {
      channels.forEach(channel => delete state[network][channel]);
    },

    [actions.CLOSE_PRIVATE_CHAT](state, { network, nick }) {
      delete state[network][nick];
    },

    [actions.socket.CHANNEL_FORWARD](state, { network, old }) {
      if (state[network]) {
        delete state[network][old];
      }
    },

    [actions.UPDATE_MESSAGE_HEIGHT](
      state,
      { wrapWidth, charWidth, windowWidth }
    ) {
      Object.keys(state).forEach(network =>
        Object.keys(state[network]).forEach(target =>
          state[network][target].forEach(message => {
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

    [actions.INIT](state, { networks }) {
      initNetworks(state, networks);
    },

    [actions.socket.NETWORKS](state, { data }) {
      initNetworks(state, data);
    }
  }
);

export function getMessageTab(network, to) {
  if (!to || to === '*' || (!isChannel(to) && to.indexOf('.') !== -1)) {
    return network;
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
            network: tab.network,
            channel: tab.name,
            next: first.id
          }
        }
      });
    }
  };
}

export function addFetchedMessages(network, tab) {
  return {
    type: actions.ADD_FETCHED_MESSAGES,
    network,
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

export function sendMessage(content, to, network) {
  return (dispatch, getState) => {
    const state = getState();
    const { wrapWidth, charWidth, windowWidth } = getApp(state);

    dispatch({
      type: actions.ADD_MESSAGE,
      network,
      tab: to,
      message: {
        from: state.networks[network].nick,
        content
      },
      wrapWidth,
      charWidth,
      windowWidth,
      socket: {
        type: 'message',
        data: { content, to, network }
      }
    });
  };
}

export function addMessage(message, network, to) {
  const tab = getMessageTab(network, to);

  return (dispatch, getState) => {
    const { wrapWidth, charWidth, windowWidth } = getApp(getState());

    dispatch({
      type: actions.ADD_MESSAGE,
      network,
      tab,
      message,
      wrapWidth,
      charWidth,
      windowWidth
    });
  };
}

export function addMessages(messages, network, to, prepend, next) {
  const tab = getMessageTab(network, to);

  return (dispatch, getState) => {
    const state = getState();

    if (next) {
      messages[0].id = next;
      messages[0].next = true;
    }

    const { wrapWidth, charWidth, windowWidth } = getApp(state);

    dispatch({
      type: actions.ADD_MESSAGES,
      network,
      tab,
      messages,
      prepend,
      wrapWidth,
      charWidth,
      windowWidth
    });
  };
}

export function addEvent(network, tab, type, ...params) {
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
    network,
    tab
  );
}

export function broadcastEvent(network, channels, type, ...params) {
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
    network
  );
}

export function broadcast(message, network, channels) {
  return addMessages(
    channels.map(channel => ({
      tab: channel,
      content: message,
      type: 'info'
    })),
    network
  );
}

export function print(message, network, channel, type) {
  if (Array.isArray(message)) {
    return addMessages(
      message.map(line => ({
        content: line,
        type
      })),
      network,
      channel
    );
  }

  return addMessage(
    {
      content: message,
      type
    },
    network,
    channel
  );
}

export function inform(message, network, channel) {
  return print(message, network, channel, 'info');
}

export function runCommand(command, channel, network) {
  return {
    type: actions.COMMAND,
    command,
    channel,
    network
  };
}

export function raw(message, network) {
  return {
    type: actions.RAW,
    message,
    network,
    socket: {
      type: 'raw',
      data: { message, network }
    }
  };
}
