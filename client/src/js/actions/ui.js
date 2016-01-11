import * as actions from '../actions';

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
