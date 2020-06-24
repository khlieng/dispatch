import reducer, { connect, setNetworkName } from '../networks';
import * as actions from '../actions';

describe('network reducer', () => {
  it('adds the network on CONNECT', () => {
    let state = reducer(
      undefined,
      connect({ host: '127.0.0.1', nick: 'nick' })
    );

    expect(state).toEqual({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null,
        connected: false,
        error: null,
        features: {}
      }
    });

    state = reducer(state, connect({ host: '127.0.0.1', nick: 'nick' }));

    expect(state).toEqual({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null,
        connected: false,
        error: null,
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
        connected: false,
        error: null,
        features: {}
      },
      '127.0.0.2': {
        name: 'srv',
        nick: 'nick',
        editedNick: null,
        connected: false,
        error: null,
        features: {}
      }
    });
  });

  it('removes the network on DISCONNECT', () => {
    let state = {
      srv: {},
      srv2: {}
    };

    state = reducer(state, {
      type: actions.DISCONNECT,
      network: 'srv2'
    });

    expect(state).toEqual({
      srv: {}
    });
  });

  it('handles SET_NETWORK_NAME', () => {
    let state = {
      srv: {
        name: 'cake'
      }
    };

    state = reducer(state, setNetworkName('pie', 'srv'));

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
      network: '127.0.0.1',
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
      network: '127.0.0.1',
      nick: 'nick2',
      editing: true
    });
    state = reducer(state, {
      type: actions.SET_NICK,
      network: '127.0.0.1',
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
      network: '127.0.0.1',
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
      network: '127.0.0.1',
      nick: 'nick2',
      editing: true
    });
    state = reducer(state, {
      type: actions.socket.NICK_FAIL,
      network: '127.0.0.1'
    });

    expect(state).toMatchObject({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null
      }
    });
  });

  it('adds the networks on INIT', () => {
    let state = reducer(undefined, {
      type: actions.INIT,
      networks: [
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
        }
      ]
    });

    expect(state).toEqual({
      '127.0.0.1': {
        name: 'stuff',
        nick: 'nick',
        editedNick: null,
        connected: true,
        features: {}
      },
      '127.0.0.2': {
        name: 'stuffz',
        nick: 'nick2',
        editedNick: null,
        connected: false,
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
      network: '127.0.0.1',
      connected: true
    });

    expect(state).toEqual({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null,
        connected: true,
        features: {}
      }
    });

    state = reducer(state, {
      type: actions.socket.CONNECTION_UPDATE,
      network: '127.0.0.1',
      connected: false,
      error: 'Bad stuff happened'
    });

    expect(state).toEqual({
      '127.0.0.1': {
        name: '127.0.0.1',
        nick: 'nick',
        editedNick: null,
        connected: false,
        error: 'Bad stuff happened',
        features: {}
      }
    });
  });
});
