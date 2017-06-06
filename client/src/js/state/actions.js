export const INVITE = 'INVITE';
export const JOIN = 'JOIN';
export const KICK = 'KICK';
export const PART = 'PART';
export const SET_TOPIC = 'SET_TOPIC';

export const SET_ENVIRONMENT = 'SET_ENVIRONMENT';

export const INPUT_HISTORY_ADD = 'INPUT_HISTORY_ADD';
export const INPUT_HISTORY_DECREMENT = 'INPUT_HISTORY_DECREMENT';
export const INPUT_HISTORY_INCREMENT = 'INPUT_HISTORY_INCREMENT';
export const INPUT_HISTORY_RESET = 'INPUT_HISTORY_RESET';

export const ADD_FETCHED_MESSAGES = 'ADD_FETCHED_MESSAGES';
export const ADD_MESSAGE = 'ADD_MESSAGE';
export const ADD_MESSAGES = 'ADD_MESSAGES';
export const COMMAND = 'COMMAND';
export const FETCH_MESSAGES = 'FETCH_MESSAGES';
export const RAW = 'RAW';
export const UPDATE_MESSAGE_HEIGHT = 'UPDATE_MESSAGE_HEIGHT';

export const CLOSE_PRIVATE_CHAT = 'CLOSE_PRIVATE_CHAT';
export const OPEN_PRIVATE_CHAT = 'OPEN_PRIVATE_CHAT';

export const SEARCH_MESSAGES = 'SEARCH_MESSAGES';
export const TOGGLE_SEARCH = 'TOGGLE_SEARCH';

export const AWAY = 'AWAY';
export const CONNECT = 'CONNECT';
export const DISCONNECT = 'DISCONNECT';
export const SET_NICK = 'SET_NICK';
export const WHOIS = 'WHOIS';

export const SET_CERT = 'SET_CERT';
export const SET_CERT_ERROR = 'SET_CERT_ERROR';
export const SET_KEY = 'SET_KEY';
export const UPLOAD_CERT = 'UPLOAD_CERT';

export const SELECT_TAB = 'SELECT_TAB';

export const HIDE_MENU = 'HIDE_MENU';
export const TOGGLE_MENU = 'TOGGLE_MENU';
export const TOGGLE_USERLIST = 'TOGGLE_USERLIST';

export function socketAction(type) {
  return `SOCKET_${type.toUpperCase()}`;
}

function createSocketActions(types) {
  const actions = {};
  types.forEach(type => {
    actions[type.toUpperCase()] = socketAction(type);
  });
  return actions;
}

export const socket = createSocketActions([
  'cert_fail',
  'cert_success',
  'channels',
  'connection_update',
  'join',
  'message',
  'mode',
  'nick',
  'part',
  'pm',
  'quit',
  'search',
  'servers',
  'topic',
  'users'
]);
