import createReducer from 'utils/createReducer';
import { LOCATION_CHANGED } from 'utils/router';
import * as actions from './actions';

const initialState = {
  showTabList: false,
  showUserList: false
};

export const getShowTabList = state => state.ui.showTabList;
export const getShowUserList = state => state.ui.showUserList;

function setMenuHidden(state) {
  state.showTabList = false;
}

export default createReducer(initialState, {
  [actions.TOGGLE_MENU](state) {
    state.showTabList = !state.showTabList;
  },

  [actions.HIDE_MENU]: setMenuHidden,
  [LOCATION_CHANGED]: setMenuHidden,

  [actions.TOGGLE_USERLIST](state) {
    state.showUserList = !state.showUserList;
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
