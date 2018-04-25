import createReducer from 'utils/createReducer';
import * as actions from './actions';

const HISTORY_MAX_LENGTH = 128;

const initialState = {
  history: [],
  index: 0
};

export const getCurrentInputHistoryEntry = state => {
  if (state.input.index === -1) {
    return null;
  }

  return state.input.history[state.input.index];
};

export default createReducer(initialState, {
  [actions.INPUT_HISTORY_ADD](state, { line }) {
    if (line.trim() && line !== state.history[0]) {
      if (history.length === HISTORY_MAX_LENGTH) {
        state.history.pop();
      }
      state.history.unshift(line);
    }
  },

  [actions.INPUT_HISTORY_RESET](state) {
    state.index = -1;
  },

  [actions.INPUT_HISTORY_INCREMENT](state) {
    if (state.index < state.history.length - 1) {
      state.index++;
    }
  },

  [actions.INPUT_HISTORY_DECREMENT](state) {
    if (state.index >= 0) {
      state.index--;
    }
  }
});

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
