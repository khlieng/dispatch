export default function createCommandMiddleware(type, handlers) {
  return ({ dispatch, getState }) => next => action => {
    if (action.type === type) {
      const words = action.command.slice(1).split(' ');
      const command = words[0];
      const params = words.slice(1);

      if (command in handlers) {
        handlers[command]({
          dispatch,
          getState,
          server: action.server,
          channel: action.channel
        }, ...params);
      }
    }

    return next(action);
  };
}
