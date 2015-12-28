import { Map } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from '../actions';

export default createReducer(Map(), {
  [actions.SET_ENVIRONMENT](state, action) {
    return state.set(action.key, action.value);
  }
});
