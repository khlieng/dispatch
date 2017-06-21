import { Record } from 'immutable';
import createReducer from 'util/createReducer';
import { LOCATION_CHANGED } from 'util/router';
import * as actions from './actions';

const State = Record({
  showTabList: false,
  showUserList: false
});

export const getShowTabList = state => state.ui.showTabList;
export const getShowUserList = state => state.ui.showUserList;

function setMenuHidden(state) {
  return state.set('showTabList', false);
}

export default createReducer(new State(), {
  [actions.TOGGLE_MENU](state) {
    return state.update('showTabList', show => !show);
  },

  [actions.HIDE_MENU]: setMenuHidden,
  [LOCATION_CHANGED]: setMenuHidden,

  [actions.TOGGLE_USERLIST](state) {
    return state.update('showUserList', show => !show);
  }
});

export function hideMenu() {
  return {
    type: actions.HIDE_MENU
  };
}

export function toggleMenu() {
  return {
    type: actions.TOGGLE_MENU
  };
}

export function toggleUserList() {
  return {
    type: actions.TOGGLE_USERLIST
  };
}
