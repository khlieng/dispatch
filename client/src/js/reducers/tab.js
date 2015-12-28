import { Record, List } from 'immutable';
import { UPDATE_PATH } from 'redux-simple-router';
import createReducer from '../util/createReducer';
import * as actions from '../actions';

const Tab = Record({
  server: null,
  channel: null,
  user: null
});

const State = Record({
  selected: new Tab(),
  history: List()
});

export default createReducer(new State(), {
  [actions.SELECT_TAB](state, action) {
    const tab = new Tab(action);
    return state
      .set('selected', tab)
      .update('history', history => history.push(tab));
  },

  [actions.PART](state, action) {
    return state.set('history', state.history.filter(tab =>
      !(tab.server === action.server && tab.channel === action.channels[0])
    ));
  },

  [actions.CLOSE_PRIVATE_CHAT](state, action) {
    return state.set('history', state.history.filter(tab =>
      !(tab.server === action.server && tab.user === action.nick)
    ));
  },

  [actions.DISCONNECT](state, action) {
    return state.set('history', state.history.filter(tab => tab.server !== action.server));
  },

  [UPDATE_PATH](state, action) {
    if (action.payload.path.indexOf('.') === -1 && state.selected.server) {
      return state.set('selected', new Tab());
    }

    return state;
  }
});
