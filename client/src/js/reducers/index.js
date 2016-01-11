import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router';
import channels from './channels';
import environment from './environment';
import input from './input';
import messages from './messages';
import privateChats from './privateChats';
import search from './search';
import servers from './servers';
import settings from './settings';
import showMenu from './showMenu';
import tab from './tab';

export default combineReducers({
  routing: routeReducer,
  channels,
  environment,
  input,
  messages,
  privateChats,
  search,
  servers,
  settings,
  showMenu,
  tab
});
