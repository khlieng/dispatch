import Immutable from 'immutable';
import reducer from '../reducers/servers';
import * as actions from '../actions';
import { connect } from '../actions/server';

describe('reducers/servers', () => {
  it('adds the server on CONNECT', () => {
    let state = reducer(undefined, connect('127.0.0.1:1337', 'nick', {}));

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        connected: false,
        name: '127.0.0.1',
        nick: 'nick'
      }
    });

    state = reducer(state, connect('127.0.0.1:1337', 'nick', {}));

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        connected: false,
        name: '127.0.0.1',
        nick: 'nick'
      }
    });

    state = reducer(state, connect('127.0.0.2:1337', 'nick', {
      name: 'srv'
    }));

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        connected: false,
        name: '127.0.0.1',
        nick: 'nick'
      },
      '127.0.0.2': {
        connected: false,
        name: 'srv',
        nick: 'nick'
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

  it('updates the nick on SOCKET_NICK', () => {
    let state = reducer(undefined, connect('127.0.0.1:1337', 'nick', {}));
    state = reducer(state, {
      type: actions.SOCKET_NICK,
      server: '127.0.0.1',
      old: 'nick',
      new: 'nick2'
    });

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        connected: false,
        name: '127.0.0.1',
        nick: 'nick2'
      }
    });
  });

  it('adds the servers on SOCKET_SERVERS', () => {
    let state = reducer(undefined, {
      type: 'SOCKET_SERVERS',
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
        connected: true
      },
      '127.0.0.2': {
        name: 'stuffz',
        nick: 'nick2',
        connected: false
      }
    });
  });

  it('updates connection status on SOCKET_CONNECTION_UPDATE', () => {
    let state = reducer(undefined, connect('127.0.0.1:1337', 'nick', {}));
    state = reducer(state, {
      type: actions.SOCKET_CONNECTION_UPDATE,
      '127.0.0.1': true
    });

    expect(state.toJS()).toEqual({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        connected: true
      }
    });
  });
});
