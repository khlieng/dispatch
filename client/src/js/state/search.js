import { List, Record } from 'immutable';
import createReducer from 'util/createReducer';
import * as actions from './actions';

const State = Record({
  show: false,
  results: List()
});

export const getSearch = state => state.search;

export default createReducer(new State(), {
  [actions.socket.SEARCH](state, action) {
    return state.set('results', List(action.results));
  },

  [actions.TOGGLE_SEARCH](state) {
    return state.set('show', !state.show);
  }
});

export function searchMessages(server, channel, phrase) {
  return {
    type: actions.SEARCH_MESSAGES,
    server,
    channel,
    phrase,
    socket: {
      type: 'search',
      data: { server, channel, phrase }
    }
  };
}

export function toggleSearch() {
  return {
    type: actions.TOGGLE_SEARCH
  };
}
