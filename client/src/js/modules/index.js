import documentTitle from './documentTitle';
import handleSocket from './handleSocket';
import initialState from './initialState';
import storage from './storage';

export default function runModules(ctx) {
  initialState(ctx);

  documentTitle(ctx);
  handleSocket(ctx);
  storage(ctx);
}
