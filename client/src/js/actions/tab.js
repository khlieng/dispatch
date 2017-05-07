import * as actions from '../actions';
import { push, replace } from '../util/router';

export function select(server, name) {
  if (name) {
    return push(`/${server}/${encodeURIComponent(name)}`);
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
      dispatch(select(tab.server, tab.name));
    } else if (servers.has(server)) {
      dispatch(select(server));
    } else {
      dispatch(replace('/'));
    }
  };
}

export function setSelectedTab(server, name = null) {
  return {
    type: actions.SELECT_TAB,
    server,
    name
  };
}
