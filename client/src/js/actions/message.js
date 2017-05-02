import * as actions from '../actions';
import { findBreakpoints, messageHeight, linkify, timestamp } from '../util';
import { getSelectedMessages } from '../reducers/messages';

let nextID = 0;

function initMessage(message, server, tab, state) {
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
    const from = message.from;
    message.from = null;
    message.type = 'action';
    message.content = from + message.content.slice(7, -1);
  }

  const charWidth = state.environment.get('charWidth');
  const wrapWidth = state.environment.get('wrapWidth');

  message.length = message.content.length;
  message.breakpoints = findBreakpoints(message.content);
  message.height = messageHeight(message, wrapWidth, charWidth, 6 * charWidth);
  message.content = linkify(message.content);

  return message;
}

function getMessageTab(server, to) {
  if (!to || to.indexOf('.') !== -1) {
    return server;
  }
  return to;
}

export function fetchMessages() {
  return (dispatch, getState) => {
    const state = getState();
    const first = getSelectedMessages(state).get(0);

    if (!first) {
      return;
    }

    const tab = state.tab.selected;
    if (tab.isChannel()) {
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

export function updateMessageHeight() {
  return (dispatch, getState) => dispatch({
    type: actions.UPDATE_MESSAGE_HEIGHT,
    wrapWidth: getState().environment.get('wrapWidth'),
    charWidth: getState().environment.get('charWidth')
  });
}

export function sendMessage(content, to, server) {
  return (dispatch, getState) => {
    const state = getState();

    dispatch({
      type: actions.SEND_MESSAGE,
      server,
      tab: to,
      message: initMessage({
        from: state.servers.getIn([server, 'nick']),
        content
      }, server, to, state),
      socket: {
        type: 'message',
        data: { content, to, server }
      }
    });
  };
}

export function addMessage(message, server, to) {
  const tab = getMessageTab(server, to);

  return (dispatch, getState) => dispatch({
    type: actions.ADD_MESSAGE,
    server,
    tab,
    message: initMessage(message, server, tab, getState())
  });
}

export function addMessages(messages, server, to, prepend, next) {
  const tab = getMessageTab(server, to);

  return (dispatch, getState) => {
    const state = getState();

    if (next) {
      messages[0].id = next;
    }

    messages.forEach(message => initMessage(message, server, message.tab || tab, state));

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
  return addMessages(channels.map(channel => ({
    tab: channel,
    content: message,
    type: 'info'
  })), server);
}

export function inform(message, server, channel) {
  if (Array.isArray(message)) {
    return addMessages(message.map(line => ({
      content: line,
      type: 'info'
    })), server, channel);
  }

  return addMessage({
    content: message,
    type: 'info'
  }, server, channel);
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
