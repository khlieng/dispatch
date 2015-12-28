import * as actions from '../actions';

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
