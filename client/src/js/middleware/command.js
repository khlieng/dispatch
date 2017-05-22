import { inform } from '../actions/message';

const notFound = 'commandNotFound';

function createContext({ dispatch, getState }, { server, channel }) {
  return { dispatch, getState, server, channel };
}

export default function createCommandMiddleware(type, handlers) {
  return store => next => action => {
    if (action.type === type) {
      const words = action.command.slice(1).split(' ');
      const command = words[0];
      const params = words.slice(1);

      let result;

      if (command in handlers) {
        result = handlers[command](createContext(store, action), ...params);
      } else if (notFound in handlers) {
        result = handlers[notFound](createContext(store, action), command);
      }

      if (typeof result === 'string' || Array.isArray(result)) {
        store.dispatch(inform(result, action.server, action.channel));
      }
    }

    return next(action);
  };
}
