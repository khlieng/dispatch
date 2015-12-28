import { List, Record } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from '../actions';

const State = Record({
  show: false,
  results: List()
});

export default createReducer(new State(), {
  [actions.SOCKET_SEARCH](state, action) {
    return state.set('results', List(action.results));
  },

  [actions.TOGGLE_SEARCH](state) {
    return state.set('show', !state.show);
  }
});
