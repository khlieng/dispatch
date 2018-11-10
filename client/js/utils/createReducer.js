import produce from 'immer';
import has from 'lodash/has';

export default function createReducer(initialState, handlers) {
  return function reducer(state = initialState, action) {
    if (has(handlers, action.type)) {
      return produce(state, draft => handlers[action.type](draft, action));
    }
    return state;
  };
}
