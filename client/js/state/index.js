import { combineReducers } from 'redux';
import app from './app';
import channels from './channels';
import input from './input';
import messages from './messages';
import privateChats from './privateChats';
import search from './search';
import servers from './servers';
import settings from './settings';
import tab from './tab';
import ui from './ui';

export * from './selectors';
export const getRouter = state => state.router;

export default function createReducer(router) {
  return combineReducers({
    router,
    app,
    channels,
    input,
    messages,
    privateChats,
    search,
    servers,
    settings,
    tab,
    ui
  });
}
