import { Record, List } from 'immutable';
import { LOCATION_CHANGE } from 'react-router-redux';
import createReducer from '../util/createReducer';
import * as actions from '../actions';

const TabRecord = Record({
  server: null,
  name: null
});

class Tab extends TabRecord {
  isChannel() {
    return this.name && this.name.charAt(0) === '#';
  }
}

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
      !(tab.server === action.server && tab.name === action.channels[0])
    ));
  },

  [actions.CLOSE_PRIVATE_CHAT](state, action) {
    return state.set('history', state.history.filter(tab =>
      !(tab.server === action.server && tab.name === action.nick)
    ));
  },

  [actions.DISCONNECT](state, action) {
    return state.set('history', state.history.filter(tab => tab.server !== action.server));
  },

  [LOCATION_CHANGE](state, action) {
    if (action.payload.pathname.indexOf('.') === -1 && state.selected.server) {
      return state.set('selected', new Tab());
    }

    return state;
  }
});
