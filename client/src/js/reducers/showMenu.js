import createReducer from '../util/createReducer';
import * as actions from '../actions';

export default createReducer(false, {
  [actions.TOGGLE_MENU](state) {
    return !state;
  },

  [actions.HIDE_MENU]() {
    return false;
  }
});
