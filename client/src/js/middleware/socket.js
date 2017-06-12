import debounce from 'lodash/debounce';

const debounceKey = action => {
  const key = action.socket.debounce.key;
  if (key) {
    return `${action.type} ${key}`;
  }
  return action.type;
};

export default function createSocketMiddleware(socket) {
  return () => next => {
    const debounced = {};

    return action => {
      if (action.socket) {
        if (action.socket.debounce) {
          const key = debounceKey(action);

          if (!debounced[key]) {
            debounced[key] = debounce((type, data) => {
              socket.send(type, data);
              debounced[key] = undefined;
            }, action.socket.debounce.delay);
          }

          debounced[key](action.socket.type, action.socket.data);
        } else {
          socket.send(action.socket.type, action.socket.data);
        }
      }

      return next(action);
    };
  };
}
