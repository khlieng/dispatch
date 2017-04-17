import * as actions from '../actions';
import { findBreakpoints, messageHeight, linkify, timestamp } from '../util';

function initMessage(message, state) {
  message.dest = message.to || message.from || message.server;
  if (message.from && message.from.indexOf('.') !== -1) {
    message.dest = message.server;
  }

  if (message.dest.charAt(0) === '#') {
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
      message: initMessage({
        from: state.servers.getIn([server, 'nick']),
        content,
        to,
        server,
        time: timestamp()
      }, state),
      socket: {
        type: 'message',
        data: { content, to, server }
      }
    });
  };
}

export function addMessage(message) {
  message.time = timestamp();

  return (dispatch, getState) => dispatch({
    type: actions.ADD_MESSAGE,
    message: initMessage(message, getState())
  });
}

export function addMessages(messages) {
  const now = timestamp();

  return (dispatch, getState) => {
    const state = getState();

    messages.forEach(message => {
      initMessage(message, state).time = now;
    });

    dispatch({
      type: actions.ADD_MESSAGES,
      messages
    });
  };
}

export function broadcast(message, server, channels) {
  return addMessages(channels.map(channel => ({
    server,
    to: channel,
    content: message,
    type: 'info'
  })));
}

export function inform(message, server, channel) {
  if (Array.isArray(message)) {
    return addMessages(message.map(line => ({
      server,
      to: channel,
      content: line,
      type: 'info'
    })));
  }

  return addMessage({
    server,
    to: channel,
    content: message,
    type: 'info'
  });
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
