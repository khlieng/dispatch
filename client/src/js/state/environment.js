import { Map } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from './actions';

export const getEnvironment = state => state.environment;

export const getWrapWidth = state => state.environment.get('wrapWidth');
export const getCharWidth = state => state.environment.get('charWidth');
export const getWindowWidth = state => state.environment.get('windowWidth');

export const getConnectDefaults = state => state.environment.get('connect_defaults');

const initialState = Map({
  connected: true
});

export default createReducer(initialState, {
  [actions.SET_ENVIRONMENT](state, action) {
    return state.set(action.key, action.value);
  },

  [actions.UPDATE_MESSAGE_HEIGHT](state, action) {
    return state
      .set('wrapWidth', action.wrapWidth)
      .set('charWidth', action.charWidth)
      .set('windowWidth', action.windowWidth);
  }
});

export function setEnvironment(key, value) {
  return {
    type: actions.SET_ENVIRONMENT,
    key,
    value
  };
}

export function setWrapWidth(width) {
  return setEnvironment('wrapWidth', width);
}

export function setCharWidth(width) {
  return setEnvironment('charWidth', width);
}
