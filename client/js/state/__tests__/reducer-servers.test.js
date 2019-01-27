import reducer, { connect, setServerName } from '../servers';
import * as actions from '../actions';

describe('server reducer', () => {
  it('adds the server on CONNECT', () => {
    let state = reducer(
      undefined,
      connect({ host: '127.0.0.1', nick: 'nick' })
    );

    expect(state).toEqual({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null,
        status: {
          connected: false,
          error: null
        },
        features: {}
      }
    });

    state = reducer(state, connect({ host: '127.0.0.1', nick: 'nick' }));

    expect(state).toEqual({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null,
        status: {
          connected: false,
          error: null
        },
        features: {}
      }
    });

    state = reducer(
      state,
      connect({ host: '127.0.0.2', nick: 'nick', name: 'srv' })
    );

    expect(state).toEqual({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null,
        status: {
          connected: false,
          error: null
        },
        features: {}
      },
      '127.0.0.2': {
        name: 'srv',
        nick: 'nick',
        editedNick: null,
        status: {
          connected: false,
          error: null
        },
        features: {}
      }
    });
  });

  it('removes the server on DISCONNECT', () => {
    let state = {
      srv: {},
      srv2: {}
    };

    state = reducer(state, {
      type: actions.DISCONNECT,
      server: 'srv2'
    });

    expect(state).toEqual({
      srv: {}
    });
  });

  it('handles SET_SERVER_NAME', () => {
    let state = {
      srv: {
        name: 'cake'
      }
    };

    state = reducer(state, setServerName('pie', 'srv'));

    expect(state).toEqual({
      srv: {
        name: 'pie'
      }
    });
  });

  it('sets editedNick when editing the nick', () => {
    let state = reducer(
      undefined,
      connect({ host: '127.0.0.1', nick: 'nick' })
    );
    state = reducer(state, {
      type: actions.SET_NICK,
      server: '127.0.0.1',
      nick: 'nick2',
      editing: true
    });

    expect(state).toMatchObject({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: 'nick2'
      }
    });
  });

  it('clears editedNick when receiving an empty nick after editing finishes', () => {
    let state = reducer(
      undefined,
      connect({ host: '127.0.0.1', nick: 'nick' })
    );
    state = reducer(state, {
      type: actions.SET_NICK,
      server: '127.0.0.1',
      nick: 'nick2',
      editing: true
    });
    state = reducer(state, {
      type: actions.SET_NICK,
      server: '127.0.0.1',
      nick: ''
    });

    expect(state).toMatchObject({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null
      }
    });
  });

  it('updates the nick on SOCKET_NICK', () => {
    let state = reducer(
      undefined,
      connect({ host: '127.0.0.1', nick: 'nick' })
    );
    state = reducer(state, {
      type: actions.socket.NICK,
      server: '127.0.0.1',
      oldNick: 'nick',
      newNick: 'nick2'
    });

    expect(state).toMatchObject({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick2',
        editedNick: null
      }
    });
  });

  it('clears editedNick on SOCKET_NICK_FAIL', () => {
    let state = reducer(
      undefined,
      connect({ host: '127.0.0.1', nick: 'nick' })
    );
    state = reducer(state, {
      type: actions.SET_NICK,
      server: '127.0.0.1',
      nick: 'nick2',
      editing: true
    });
    state = reducer(state, {
      type: actions.socket.NICK_FAIL,
      server: '127.0.0.1'
    });

    expect(state).toMatchObject({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null
      }
    });
  });

  it('adds the servers on SOCKET_SERVERS', () => {
    let state = reducer(undefined, {
      type: actions.socket.SERVERS,
      data: [
        {
          host: '127.0.0.1',
          name: 'stuff',
          nick: 'nick',
          status: {
            connected: true
          }
        },
        {
          host: '127.0.0.2',
          name: 'stuffz',
          nick: 'nick2',
          status: {
            connected: false
          }
        }
      ]
    });

    expect(state).toEqual({
      '127.0.0.1': {
        name: 'stuff',
        nick: 'nick',
        editedNick: null,
        status: {
          connected: true
        },
        features: {}
      },
      '127.0.0.2': {
        name: 'stuffz',
        nick: 'nick2',
        editedNick: null,
        status: {
          connected: false
        },
        features: {}
      }
    });
  });

  it('updates connection status on SOCKET_CONNECTION_UPDATE', () => {
    let state = reducer(
      undefined,
      connect({ host: '127.0.0.1', nick: 'nick' })
    );
    state = reducer(state, {
      type: actions.socket.CONNECTION_UPDATE,
      server: '127.0.0.1',
      connected: true
    });

    expect(state).toEqual({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null,
        status: {
          connected: true
        },
        features: {}
      }
    });

    state = reducer(state, {
      type: actions.socket.CONNECTION_UPDATE,
      server: '127.0.0.1',
      connected: false,
      error: 'Bad stuff happened'
    });

    expect(state).toEqual({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null,
        status: {
          connected: false,
          error: 'Bad stuff happened'
        },
        features: {}
      }
    });
  });
});
