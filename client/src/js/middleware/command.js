export default function createCommandMiddleware(type, handlers) {
  return store => next => action => {
    if (action.type === type) {
      const words = action.command.slice(1).split(' ');
      const command = words[0];
      const params = words.slice(1);

      if (Object.prototype.hasOwnProperty.call(handlers, command)) {
        handlers[command]({
          dispatch: store.dispatch,
          getState: store.getState,
          server: action.server,
          channel: action.channel
        }, ...params);
      }
    }

    return next(action);
  };
}
