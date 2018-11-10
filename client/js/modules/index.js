import documentTitle from './documentTitle';
import fonts from './fonts';
import initialState from './initialState';
import socket from './socket';
import storage from './storage';
import widthUpdates from './widthUpdates';

export default function runModules(ctx) {
  fonts(ctx);
  initialState(ctx);

  documentTitle(ctx);
  socket(ctx);
  storage(ctx);
  widthUpdates(ctx);
}
