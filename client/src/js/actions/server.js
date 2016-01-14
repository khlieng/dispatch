import * as actions from '../actions';
import { updateSelection } from './tab';

export function connect(server, nick, options) {
  let host = server;
  const i = server.indexOf(':');
  if (i > 0) {
    host = server.slice(0, i);
  }

  return {
    type: actions.CONNECT,
    host,
    nick,
    options,
    socket: {
      type: 'connect',
      data: {
        server,
        nick,
        username: options.username || nick,
        password: options.password,
        realname: options.realname || nick,
        tls: options.tls || false,
        name: options.name || server
      }
    }
  };
}

export function disconnect(server) {
  return dispatch => {
    dispatch({
      type: actions.DISCONNECT,
      server,
      socket: {
        type: 'quit',
        data: { server }
      }
    });
    dispatch(updateSelection());
  };
}

export function whois(user, server) {
  return {
    type: actions.WHOIS,
    user,
    server,
    socket: {
      type: 'whois',
      data: { user, server }
    }
  };
}

export function away(message, server) {
  return {
    type: actions.AWAY,
    message,
    server,
    socket: {
      type: 'away',
      data: { message, server }
    }
  };
}

export function setNick(nick, server) {
  return {
    type: actions.SET_NICK,
    nick,
    server,
    socket: {
      type: 'nick',
      data: {
        new: nick,
        server
      }
    }
  };
}
