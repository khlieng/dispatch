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
  message.message = message.message.replace(/\s\s+/g, ' ');

  if (message.message.indexOf('\x01ACTION') === 0) {
    const from = message.from;
    message.from = null;
    message.type = 'action';
    message.message = from + message.message.slice(7, -1);
  }

  const charWidth = state.environment.get('charWidth');
  const wrapWidth = state.environment.get('wrapWidth');

  message.length = message.message.length;
  message.breakpoints = findBreakpoints(message.message);
  message.height = messageHeight(message, wrapWidth, charWidth, 6 * charWidth);
  message.message = linkify(message.message);

  return message;
}

export function updateMessageHeight() {
  return (dispatch, getState) => dispatch({
    type: actions.UPDATE_MESSAGE_HEIGHT,
    wrapWidth: getState().environment.get('wrapWidth'),
    charWidth: getState().environment.get('charWidth')
  });
}

export function sendMessage(message, to, server) {
  return (dispatch, getState) => {
    const state = getState();

    dispatch({
      type: actions.SEND_MESSAGE,
      message: initMessage({
        from: state.servers.getIn([server, 'nick']),
        message,
        to,
        server,
        time: timestamp()
      }, state),
      socket: {
        type: 'chat',
        data: { message, to, server }
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
    message,
    type: 'info'
  })));
}

export function inform(message, server, channel) {
  if (Array.isArray(message)) {
    return addMessages(message.map(msg => ({
      server,
      to: channel,
      message: msg,
      type: 'info'
    })));
  }

  return addMessage({
    server,
    to: channel,
    message,
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
