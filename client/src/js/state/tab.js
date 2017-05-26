import { Record, List } from 'immutable';
import { push, replace, LOCATION_CHANGED } from '../util/router';
import createReducer from '../util/createReducer';
import * as actions from './actions';

const TabRecord = Record({
  server: null,
  name: null
});

class Tab extends TabRecord {
  isChannel() {
    return this.name && this.name.charAt(0) === '#';
  }

  toString() {
    let str = this.server;
    if (this.name) {
      str += `:${this.name}`;
    }
    return str;
  }
}

const State = Record({
  selected: new Tab(),
  history: List()
});

function selectTab(state, action) {
  const tab = new Tab(action);
  return state
    .set('selected', tab)
    .update('history', history => history.push(tab));
}

export const getSelectedTab = state => state.tab.selected;

export default createReducer(new State(), {
  [actions.SELECT_TAB]: selectTab,

  [actions.PART](state, action) {
    return state.set('history', state.history.filter(tab =>
      !(tab.server === action.server && tab.name === action.channels[0])
    ));
  },

  [actions.CLOSE_PRIVATE_CHAT](state, action) {
    return state.set('history', state.history.filter(tab =>
      !(tab.server === action.server && tab.name === action.nick)
    ));
  },

  [actions.DISCONNECT](state, action) {
    return state.set('history', state.history.filter(tab => tab.server !== action.server));
  },

  [LOCATION_CHANGED](state, action) {
    const { route, params } = action;
    if (route === 'chat') {
      return selectTab(state, params);
    }

    return state.set('selected', new Tab());
  }
});

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
