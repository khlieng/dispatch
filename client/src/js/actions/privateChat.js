import * as actions from '../actions';
import { updateSelection } from './tab';

export function openPrivateChat(server, nick) {
  return {
    type: actions.OPEN_PRIVATE_CHAT,
    server,
    nick
  };
}

export function closePrivateChat(server, nick) {
  return dispatch => {
    dispatch({
      type: actions.CLOSE_PRIVATE_CHAT,
      server,
      nick
    });
    dispatch(updateSelection());
  };
}
