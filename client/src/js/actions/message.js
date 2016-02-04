import * as actions from '../actions';

export function sendMessage(message, to, server) {
  return (dispatch, getState) => dispatch({
    type: actions.SEND_MESSAGE,
    from: getState().servers.getIn([server, 'nick']),
    message,
    to,
    server,
    time: new Date(),
    socket: {
      type: 'chat',
      data: { message, to, server }
    }
  });
}

export function addMessage(message) {
  message.time = new Date();

  return {
    type: actions.ADD_MESSAGE,
    message
  };
}

export function addMessages(messages) {
  const now = new Date();
  messages.forEach(message => message.time = now);

  return {
    type: actions.ADD_MESSAGES,
    messages
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
