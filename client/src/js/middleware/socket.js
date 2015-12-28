export default function createSocketMiddleware(socket) {
  return () => next => action => {
    if (action.socket) {
      socket.send(action.socket.type, action.socket.data);
    }

    return next(action);
  };
}
