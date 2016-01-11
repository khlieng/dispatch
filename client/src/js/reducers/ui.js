import { Record } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from '../actions';

const State = Record({
  showTabList: false,
  showUserList: false
});

export default createReducer(new State(), {
  [actions.TOGGLE_MENU](state) {
    return state.update('showTabList', show => !show);
  },

  [actions.HIDE_MENU](state) {
    return state.set('showTabList', false);
  },

  [actions.TOGGLE_USERLIST](state) {
    return state.update('showUserList', show => !show);
  }
});
