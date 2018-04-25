import createReducer from 'utils/createReducer';
import { push, replace, LOCATION_CHANGED } from 'utils/router';
import * as actions from './actions';

const initialState = {
  selected: {},
  history: []
};

function selectTab(state, action) {
  state.selected = {
    server: action.server,
    name: action.name
  };
  state.history.push(state.selected);
}

export const getSelectedTab = state => state.tab.selected;

export default createReducer(initialState, {
  [actions.SELECT_TAB]: selectTab,

  [actions.PART](state, action) {
    state.history = state.history.filter(
      tab => !(tab.server === action.server && tab.name === action.channels[0])
    );
  },

  [actions.CLOSE_PRIVATE_CHAT](state, action) {
    state.history = state.history.filter(
      tab => !(tab.server === action.server && tab.name === action.nick)
    );
  },

  [actions.DISCONNECT](state, action) {
    state.history = state.history.filter(tab => tab.server !== action.server);
  },

  [LOCATION_CHANGED](state, action) {
    const { route, params } = action;
    if (route === 'chat') {
      selectTab(state, params);
    } else {
      state.selected = {};
    }
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
    const { history } = state.tab;
    const { servers } = state;
    const { server } = state.tab.selected;
    const serverAddrs = Object.keys(servers);

    if (serverAddrs.length === 0) {
      dispatch(replace('/connect'));
    } else if (history.length > 0) {
      const tab = history[history.length - 1];
      dispatch(select(tab.server, tab.name, true));
    } else if (servers[server]) {
      dispatch(select(server, null, true));
    } else {
      dispatch(select(serverAddrs.sort()[0], null, true));
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
