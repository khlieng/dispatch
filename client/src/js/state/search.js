import createReducer from 'utils/createReducer';
import * as actions from './actions';

const initialState = {
  show: false,
  results: []
};

export const getSearch = state => state.search;

export default createReducer(initialState, {
  [actions.socket.SEARCH](state, { results }) {
    state.results = results;
  },

  [actions.TOGGLE_SEARCH](state) {
    state.show = !state.show;
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
