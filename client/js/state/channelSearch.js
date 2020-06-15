import createReducer from 'utils/createReducer';
import * as actions from 'state/actions';

const initialState = {
  results: [],
  end: false,
  topCache: {}
};

export default createReducer(initialState, {
  [actions.socket.CHANNEL_SEARCH](state, { results, start, network, q }) {
    if (results) {
      state.end = false;

      if (start > 0) {
        state.results.push(...results);
      } else {
        state.results = results;

        if (!q) {
          state.topCache[network] = results;
        }
      }
    } else {
      state.end = true;
    }
  },

  [actions.OPEN_MODAL](state, { name, payload }) {
    if (name === 'channel') {
      state.results = state.topCache[payload] || [];
      state.end = false;
    }
  }
});

export function searchChannels(network, q, start) {
  return {
    type: actions.CHANNEL_SEARCH,
    network,
    q,
    socket: {
      type: 'channel_search',
      data: { network, q, start }
    }
  };
}
