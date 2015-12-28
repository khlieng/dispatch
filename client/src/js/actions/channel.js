import * as actions from '../actions';
import { updateSelection } from './tab';

export function join(channels, server) {
  return {
    type: actions.JOIN,
    channels,
    server,
    socket: {
      type: 'join',
      data: { channels, server }
    }
  };
}

export function part(channels, server) {
  return dispatch => {
    dispatch({
      type: actions.PART,
      channels,
      server,
      socket: {
        type: 'part',
        data: { channels, server }
      }
    });
    dispatch(updateSelection());
  };
}

export function invite(user, channel, server) {
  return {
    type: actions.INVITE,
    user,
    channel,
    server,
    socket: {
      type: 'invite',
      data: { user, channel, server }
    }
  };
}

export function kick(user, channel, server) {
  return {
    type: actions.KICK,
    user,
    channel,
    server,
    socket: {
      type: 'kick',
      data: { user, channel, server }
    }
  };
}
