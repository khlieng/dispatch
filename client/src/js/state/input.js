import { List, Record } from 'immutable';
import createReducer from 'utils/createReducer';
import * as actions from './actions';

const HISTORY_MAX_LENGTH = 128;

const State = Record({
  history: List(),
  index: 0
});

export const getCurrentInputHistoryEntry = state => {
  if (state.input.index === -1) {
    return null;
  }

  return state.input.history.get(state.input.index);
};

export default createReducer(new State(), {
  [actions.INPUT_HISTORY_ADD](state, action) {
    const { line } = action;
    if (line.trim() && line !== state.history.get(0)) {
      if (history.length === HISTORY_MAX_LENGTH) {
        return state.set('history', state.history.unshift(line).pop());
      }

      return state.set('history', state.history.unshift(line));
    }

    return state;
  },

  [actions.INPUT_HISTORY_RESET](state) {
    return state.set('index', -1);
  },

  [actions.INPUT_HISTORY_INCREMENT](state) {
    if (state.index < state.history.size - 1) {
      return state.set('index', state.index + 1);
    }

    return state;
  },

  [actions.INPUT_HISTORY_DECREMENT](state) {
    if (state.index >= 0) {
      return state.set('index', state.index - 1);
    }

    return state;
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
