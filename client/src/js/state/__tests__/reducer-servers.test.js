import Immutable from 'immutable';
import reducer, { connect, setServerName } from '../servers';
import * as actions from '../actions';

describe('server reducer', () => {
  it('adds the server on CONNECT', () => {
    let state = reducer(undefined, connect('127.0.0.1:1337', 'nick', {}));

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        connected: false,
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null
      }
    });

    state = reducer(state, connect('127.0.0.1:1337', 'nick', {}));

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        connected: false,
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null
      }
    });

    state = reducer(state, connect('127.0.0.2:1337', 'nick', {
      name: 'srv'
    }));

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        connected: false,
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null
      },
      '127.0.0.2': {
        connected: false,
        name: 'srv',
        nick: 'nick',
        editedNick: null
      }
    });
  });

  it('removes the server on DISCONNECT', () => {
    let state = Immutable.fromJS({
      srv: {},
      srv2: {}
    });

    state = reducer(state, {
      type: actions.DISCONNECT,
      server: 'srv2'
    });

    expect(state.toJS()).toEqual({
      srv: {}
    });
  });

  it('handles SET_SERVER_NAME', () => {
    let state = Immutable.fromJS({
      srv: {
        name: 'cake'
      }
    });

    state = reducer(state, setServerName('pie', 'srv'));

    expect(state.toJS()).toEqual({
      srv: {
        name: 'pie'
      }
    });
  });

  it('sets editedNick when editing the nick', () => {
    let state = reducer(undefined, connect('127.0.0.1:1337', 'nick', {}));
    state = reducer(state, {
      type: actions.SET_NICK,
      server: '127.0.0.1',
      nick: 'nick2',
      editing: true
    });

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        connected: false,
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: 'nick2'
      }
    });
  });

  it('clears editedNick when receiving an empty nick after editing finishes', () => {
    let state = reducer(undefined, connect('127.0.0.1:1337', 'nick', {}));
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

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        connected: false,
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null
      }
    });
  });

  it('updates the nick on SOCKET_NICK', () => {
    let state = reducer(undefined, connect('127.0.0.1:1337', 'nick', {}));
    state = reducer(state, {
      type: actions.socket.NICK,
      server: '127.0.0.1',
      oldNick: 'nick',
      newNick: 'nick2'
    });

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        connected: false,
        name: '127.0.0.1',
        nick: 'nick2',
        editedNick: null
      }
    });
  });

  it('clears editedNick on SOCKET_NICK_FAIL', () => {
    let state = reducer(undefined, connect('127.0.0.1:1337', 'nick', {}));
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

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        connected: false,
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
          connected: true
        },
        {
          host: '127.0.0.2',
          name: 'stuffz',
          nick: 'nick2',
          connected: false
        },
      ]
    });

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        name: 'stuff',
        nick: 'nick',
        editedNick: null,
        connected: true
      },
      '127.0.0.2': {
        name: 'stuffz',
        nick: 'nick2',
        editedNick: null,
        connected: false
      }
    });
  });

  it('updates connection status on SOCKET_CONNECTION_UPDATE', () => {
    let state = reducer(undefined, connect('127.0.0.1:1337', 'nick', {}));
    state = reducer(state, {
      type: actions.socket.CONNECTION_UPDATE,
      '127.0.0.1': true
    });

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null,
        connected: true
      }
    });
  });
});
