import { addMessages, inform, print } from '../state/messages';
import { isChannel } from '../util';

export const beforeHandler = '_before';
export const notFoundHandler = 'commandNotFound';

function createContext({ dispatch, getState }, { server, channel }) {
  return { dispatch, getState, server, channel, isChannel: isChannel(channel) };
}

// TODO: Pull this out as convenience action
function process({ dispatch, server, channel }, result) {
  if (typeof result === 'string') {
    dispatch(inform(result, server, channel));
  } else if (Array.isArray(result)) {
    if (typeof result[0] === 'string') {
      dispatch(inform(result, server, channel));
    } else if (typeof result[0] === 'object') {
      dispatch(addMessages(result, server, channel));
    }
  } else if (typeof result === 'object' && result) {
    dispatch(print(result.content, server, channel, result.type));
  }
}

export default function createCommandMiddleware(type, handlers) {
  return store => next => action => {
    if (action.type === type) {
      const words = action.command.slice(1).split(' ');
      const command = words[0];
      const params = words.slice(1);

      if (command in handlers) {
        const ctx = createContext(store, action);
        if (beforeHandler in handlers) {
          process(ctx, handlers[beforeHandler](ctx, command, ...params));
        }
        process(ctx, handlers[command](ctx, ...params));
      } else if (notFoundHandler in handlers) {
        const ctx = createContext(store, action);
        process(ctx, handlers[notFoundHandler](ctx, command, ...params));
      }
    }

    return next(action);
  };
}
