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
    network: action.network,
    name: action.name
  };
  state.history.push(state.selected);
}

export const getSelectedTab = state => state.tab.selected;

export default createReducer(initialState, {
  [actions.SELECT_TAB]: selectTab,

  [actions.PART](state, action) {
    state.history = state.history.filter(
      tab =>
        !(tab.network === action.network && tab.name === action.channels[0])
    );
  },

  [actions.CLOSE_PRIVATE_CHAT](state, action) {
    state.history = state.history.filter(
      tab => !(tab.network === action.network && tab.name === action.nick)
    );
  },

  [actions.DISCONNECT](state, action) {
    state.history = state.history.filter(tab => tab.network !== action.network);
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

export function select(network, name, doReplace) {
  const navigate = doReplace ? replace : push;
  if (name) {
    return navigate(`/${network}/${encodeURIComponent(name)}`);
  }
  return navigate(`/${network}`);
}

export function tabExists(
  { network, name },
  { networks, channels, privateChats }
) {
  return (
    (name && get(channels, [network, name])) ||
    (!name && network && networks[network]) ||
    (name && find(privateChats[network], nick => nick === name))
  );
}

function parseTabCookie() {
  const cookie = Cookie.get('tab');
  if (cookie) {
    const [network, name = null] = cookie.split(/;(.+)/);
    return { network, name };
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
        return dispatch(select(tab.network, tab.name, true));
      }
    }

    const { networks } = state;
    const { history } = state.tab;
    const { network } = state.tab.selected;
    const networkAddrs = Object.keys(networks);

    if (networkAddrs.length === 0) {
      dispatch(replace('/connect'));
    } else if (
      history.length > 0 &&
      tabExists(history[history.length - 1], state)
    ) {
      const tab = history[history.length - 1];
      dispatch(select(tab.network, tab.name, true));
    } else if (networks[network]) {
      dispatch(select(network, null, true));
    } else {
      dispatch(select(networkAddrs.sort()[0], null, true));
    }
  };
}

export function setSelectedTab(network, name = null) {
  return {
    type: actions.SELECT_TAB,
    network,
    name
  };
}
