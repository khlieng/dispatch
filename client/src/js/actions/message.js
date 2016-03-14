import * as actions from '../actions';
import { messageHeight } from '../util';

function initMessage(message, state) {
  message.dest = message.to || message.from || message.server;
  if (message.from && message.from.indexOf('.') !== -1) {
    message.dest = message.server;
  }

  if (message.dest.charAt(0) === '#') {
    message.channel = true;
  }

  // Combine multiple adjacent spaces into a single one
  message.message = message.message.replace(/\s\s+/g, ' ');

  const charWidth = state.environment.get('charWidth');
  const wrapWidth = state.environment.get('wrapWidth');

  message.height = messageHeight(message, wrapWidth, charWidth, 6 * charWidth);

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

    dispatch(initMessage({
      type: actions.SEND_MESSAGE,
      from: state.servers.getIn([server, 'nick']),
      message,
      to,
      server,
      time: new Date(),
      socket: {
        type: 'chat',
        data: { message, to, server }
      }
    }, state));
  };
}

export function addMessage(message) {
  message.time = new Date();

  return (dispatch, getState) => dispatch({
    type: actions.ADD_MESSAGE,
    message: initMessage(message, getState())
  });
}

export function addMessages(messages) {
  const now = new Date();

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
