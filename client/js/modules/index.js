import documentTitle from './documentTitle';
import fonts from './fonts';
import initialState from './initialState';
import route from './route';
import socket from './socket';
import storage from './storage';
import widthUpdates from './widthUpdates';

export default function runModules(ctx) {
  fonts(ctx);
  initialState(ctx);
  route(ctx);

  documentTitle(ctx);
  socket(ctx);
  storage(ctx);
  widthUpdates(ctx);
}
