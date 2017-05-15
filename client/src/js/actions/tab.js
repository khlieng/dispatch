import * as actions from '../actions';
import { push, replace } from '../util/router';

export function select(server, name, doReplace) {
  const navigate = doReplace ? replace : push;
  if (name) {
    return navigate(`/${server}/${encodeURIComponent(name)}`);
  }
  return navigate(`/${server}`);
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
      dispatch(select(tab.server, tab.name, true));
    } else if (servers.has(server)) {
      dispatch(select(server, null, true));
    } else {
      dispatch(select(servers.keySeq().first(), null, true));
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
