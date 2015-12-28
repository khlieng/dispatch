import * as actions from '../actions';

export function searchMessages(server, channel, phrase) {
  return {
    type: actions.SEARCH_MESSAGES,
    server,
    channel,
    phrase,
    socket: {
      type: 'search',
      data: { server, channel, phrase }
    }
  };
}

export function toggleSearch() {
  return {
    type: actions.TOGGLE_SEARCH
  };
}
