import { ADD_MESSAGES, ADD_FETCHED_MESSAGES } from 'state/actions';

//
// This middleware handles waiting until MessageBox
// is ready before adding messages to the top
//
const message = store => next => {
  const ready = {};
  const cache = {};

  return action => {
    if (action.type === ADD_MESSAGES && action.prepend) {
      const key = `${action.server} ${action.channel}`;

      if (ready[key]) {
        ready[key] = false;
        return next(action);
      }

      cache[key] = action;
    } else if (action.type === ADD_FETCHED_MESSAGES) {
      const key = `${action.server} ${action.channel}`;
      ready[key] = true;

      if (cache[key]) {
        store.dispatch(cache[key]);
        cache[key] = undefined;
      }
    } else {
      return next(action);
    }
  };
};

export default message;
