import { push, replace } from 'react-router-redux';
import * as actions from '../actions';

export function select(server, channel, pm) {
  if (pm) {
    return push(`/${server}/pm/${channel}`);
  } else if (channel) {
    return push(`/${server}/${encodeURIComponent(channel)}`);
  }

  return push(`/${server}`);
}

export function updateSelection() {
  return (dispatch, getState) => {
    const state = getState();
    const history = state.tab.history;
    const { servers } = state;
    const { server } = state.tab.selected;

    if (servers.size === 0) {
      dispatch(replace('/connect'));
    } else if (history.size > 0) {
      const tab = history.last();
      dispatch(select(tab.server, tab.channel || tab.user, tab.user));
    } else if (servers.has(server)) {
      dispatch(select(server));
    } else {
      dispatch(replace('/'));
    }
  };
}

export function setSelectedChannel(server, channel = null) {
  return {
    type: actions.SELECT_TAB,
    server,
    channel
  };
}

export function setSelectedUser(server, user = null) {
  return {
    type: actions.SELECT_TAB,
    server,
    user
  };
}
