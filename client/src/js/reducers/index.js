import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux';
import channels from './channels';
import environment from './environment';
import input from './input';
import messages from './messages';
import privateChats from './privateChats';
import search from './search';
import servers from './servers';
import settings from './settings';
import tab from './tab';
import ui from './ui';

export default combineReducers({
  routing: routerReducer,
  channels,
  environment,
  input,
  messages,
  privateChats,
  search,
  servers,
  settings,
  tab,
  ui
});
