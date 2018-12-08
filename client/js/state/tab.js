import get from 'lodash/get';
import Cookie from 'js-cookie';
import createReducer from 'utils/createReducer';
import { push, replace, LOCATION_CHANGED } from 'utils/router';
import * as actions from './actions';
import { find } from '../utils';

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

export function tabExists(
  { server, name },
  { servers, channels, privateChats }
) {
  return (
    (name && get(channels, [server, name])) ||
    (!name && server && servers[server]) ||
    (name && find(privateChats[server], nick => nick === name))
  );
}

function parseTabCookie() {
  const cookie = Cookie.get('tab');
  if (cookie) {
    const [server, name = null] = cookie.split(/;(.+)/);
    return { server, name };
  }
  return null;
}

export function updateSelection(tryCookie) {
  return (dispatch, getState) => {
    const state = getState();

    if (tabExists(state.tab.selected, state)) {
      return;
    }

    if (tryCookie) {
      const tab = parseTabCookie();
      if (tab && tabExists(tab, state)) {
        return dispatch(select(tab.server, tab.name, true));
      }
    }

    const { servers } = state;
    const { history } = state.tab;
    const { server } = state.tab.selected;
    const serverAddrs = Object.keys(servers);

    if (serverAddrs.length === 0) {
      dispatch(replace('/connect'));
    } else if (
      history.length > 0 &&
      tabExists(history[history.length - 1], state)
    ) {
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
