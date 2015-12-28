import * as actions from '../actions';

export function addInputHistory(line) {
  return {
    type: actions.INPUT_HISTORY_ADD,
    line
  };
}

export function resetInputHistory() {
  return {
    type: actions.INPUT_HISTORY_RESET
  };
}

export function incrementInputHistory() {
  return {
    type: actions.INPUT_HISTORY_INCREMENT
  };
}

export function decrementInputHistory() {
  return {
    type: actions.INPUT_HISTORY_DECREMENT
  };
}
