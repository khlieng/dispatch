import createReducer from 'utils/createReducer';
import * as actions from 'state/actions';

const initialState = {
  results: [],
  end: false
};

export default createReducer(initialState, {
  [actions.socket.CHANNEL_SEARCH](state, { results, start }) {
    if (results) {
      state.end = false;

      if (start > 0) {
        state.results.push(...results);
      } else {
        state.results = results;
      }
    } else {
      state.end = true;
    }
  },

  [actions.OPEN_MODAL](state, { name }) {
    if (name === 'channel') {
      return initialState;
    }
  }
});

export function searchChannels(server, q, start) {
  return {
    type: actions.CHANNEL_SEARCH,
    server,
    q,
    socket: {
      type: 'channel_search',
      data: { server, q, start }
    }
  };
}
