import sortBy from 'lodash/sortBy';
import { isDM } from 'utils';
import createReducer from 'utils/createReducer';
import { updateSelection } from './tab';
import * as actions from './actions';

export const getPrivateChats = state => state.privateChats;

function open(state, network, nick) {
  if (!state[network]) {
    state[network] = [];
  }
  if (!state[network].includes(nick)) {
    state[network].push(nick);
    state[network] = sortBy(state[network], v => v.toLowerCase());
  }
}

export default createReducer(
  {},
  {
    [actions.OPEN_PRIVATE_CHAT](state, action) {
      open(state, action.network, action.nick);
    },

    [actions.CLOSE_PRIVATE_CHAT](state, { network, nick }) {
      const i = state[network]?.findIndex(n => n === nick);
      if (i !== -1) {
        state[network].splice(i, 1);
      }
    },

    [actions.INIT](state, { openDMs }) {
      if (openDMs) {
        openDMs.forEach(({ network, name }) => {
          if (!state[network]) {
            state[network] = [];
          }

          state[network].push(name);
        });
      }
    },

    [actions.ADD_MESSAGE](state, { message }) {
      if (isDM(message)) {
        open(state, message.network, message.from);
      }
    },

    [actions.DISCONNECT](state, { network }) {
      delete state[network];
    }
  }
);

export function openPrivateChat(network, nick) {
  return (dispatch, getState) => {
    if (!getState().privateChats[network]?.includes(nick)) {
      dispatch({
        type: actions.OPEN_PRIVATE_CHAT,
        network,
        nick,
        socket: {
          type: 'open_dm',
          data: { network, name: nick }
        }
      });
    }
  };
}

export function closePrivateChat(network, nick) {
  return dispatch => {
    dispatch({
      type: actions.CLOSE_PRIVATE_CHAT,
      network,
      nick,
      socket: {
        type: 'close_dm',
        data: { network, name: nick }
      }
    });
    dispatch(updateSelection());
  };
}
