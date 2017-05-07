import { Record } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from '../actions';
import { LOCATION_CHANGED } from '../util/router';

const State = Record({
  showTabList: false,
  showUserList: false
});

function hideMenu(state) {
  return state.set('showTabList', false);
}

export default createReducer(new State(), {
  [actions.TOGGLE_MENU](state) {
    return state.update('showTabList', show => !show);
  },

  [actions.HIDE_MENU]: hideMenu,
  [LOCATION_CHANGED]: hideMenu,

  [actions.TOGGLE_USERLIST](state) {
    return state.update('showUserList', show => !show);
  }
});
