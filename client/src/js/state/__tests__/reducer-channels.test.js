import Immutable from 'immutable';
import reducer from '../channels';
import { connect } from '../servers';
import * as actions from '../actions';

describe('reducers/channels', () => {
  it('removes channels on PART', () => {
    let state = Immutable.fromJS({
      srv1: {
        chan1: {}, chan2: {}, chan3: {}
      },
      srv2: {
        chan1: {}
      }
    });

    state = reducer(state, {
      type: actions.PART,
      server: 'srv1',
      channels: ['chan1', 'chan3']
    });

    expect(state.toJS()).toEqual({
      srv1: {
        chan2: {}
      },
      srv2: {
        chan1: {}
      }
    });
  });

  it('handles SOCKET_PART', () => {
    let state = reducer(undefined, socket_join('srv', 'chan1', 'nick1'));
    state = reducer(state, socket_join('srv', 'chan1', 'nick2'));
    state = reducer(state, socket_join('srv', 'chan2', 'nick2'));

    state = reducer(state, {
      type: actions.socket.PART,
      server: 'srv',
      channel: 'chan1',
      user: 'nick2'
    });

    expect(state.toJS()).toEqual({
      srv: {
        chan1: {
          users: [
            { mode: '', nick: 'nick1', renderName: 'nick1' },
          ]
        },
        chan2: {
          users: [
            { mode: '', nick: 'nick2', renderName: 'nick2' }
          ]
        }
      }
    });
  });

  it('handles SOCKET_JOIN', () => {
    const state = reducer(undefined, socket_join('srv', 'chan1', 'nick1'));

    expect(state.toJS()).toEqual({
      srv: {
        chan1: {
          users: [
            { mode: '', nick: 'nick1', renderName: 'nick1' }
          ]
        }
      }
    });
  });

  it('handles SOCKET_QUIT', () => {
    let state = reducer(undefined, socket_join('srv', 'chan1', 'nick1'));
    state = reducer(state, socket_join('srv', 'chan1', 'nick2'));
    state = reducer(state, socket_join('srv', 'chan2', 'nick2'));

    state = reducer(state, {
      type: actions.socket.QUIT,
      server: 'srv',
      user: 'nick2'
    });

    expect(state.toJS()).toEqual({
      srv: {
        chan1: {
          users: [
            { mode: '', nick: 'nick1', renderName: 'nick1' }
          ]
        },
        chan2: {
          users: []
        }
      }
    });
  });

  it('handles SOCKET_NICK', () => {
    let state = reducer(undefined, socket_join('srv', 'chan1', 'nick1'));
    state = reducer(state, socket_join('srv', 'chan1', 'nick2'));
    state = reducer(state, socket_join('srv', 'chan2', 'nick2'));

    state = reducer(state, {
      type: actions.socket.NICK,
      server: 'srv',
      old: 'nick1',
      new: 'nick3'
    });

    expect(state.toJS()).toEqual({
      srv: {
        chan1: {
          users: [
            { mode: '', nick: 'nick2', renderName: 'nick2' },
            { mode: '', nick: 'nick3', renderName: 'nick3' }
          ]
        },
        chan2: {
          users: [
            { mode: '', nick: 'nick2', renderName: 'nick2' }
          ]
        }
      }
    });
  });

  it('handles SOCKET_USERS', () => {
    const state = reducer(undefined, {
      type: actions.socket.USERS,
      server: 'srv',
      channel: 'chan1',
      users: [
        'user3',
        'user2',
        '@user4',
        'user1',
        '+user5'
      ]
    });

    expect(state.toJS()).toEqual({
      srv: {
        chan1: {
          users: [
            { mode: 'o', nick: 'user4', renderName: '@user4' },
            { mode: 'v', nick: 'user5', renderName: '+user5' },
            { mode: '', nick: 'user1', renderName: 'user1' },
            { mode: '', nick: 'user2', renderName: 'user2' },
            { mode: '', nick: 'user3', renderName: 'user3' }
          ]
        }
      }
    })
  });

  it('handles SOCKET_TOPIC', () => {
    const state = reducer(undefined, {
      type: actions.socket.TOPIC,
      server: 'srv',
      channel: 'chan1',
      topic: 'the topic'
    });

    expect(state.toJS()).toEqual({
      srv: {
        chan1: {
          topic: 'the topic'
        }
      }
    });
  });

  it('handles SOCKET_MODE', () => {
    let state = reducer(undefined, socket_join('srv', 'chan1', 'nick1'));
    state = reducer(state, socket_join('srv', 'chan1', 'nick2'));
    state = reducer(state, socket_join('srv', 'chan2', 'nick2'));

    state = reducer(state, socket_mode('srv', 'chan1', 'nick1', 'o', ''));

    expect(state.toJS()).toEqual({
      srv: {
        chan1: {
          users: [
            { mode: 'o', nick: 'nick1', renderName: '@nick1' },
            { mode: '', nick: 'nick2', renderName: 'nick2' }
          ]
        },
        chan2: {
          users: [
            { mode: '', nick: 'nick2', renderName: 'nick2' }
          ]
        }
      }
    });

    state = reducer(state, socket_mode('srv', 'chan1', 'nick1' ,'v', 'o'));
    state = reducer(state, socket_mode('srv', 'chan1', 'nick2', 'o', ''));
    state = reducer(state, socket_mode('srv', 'chan2', 'not_there', 'x', ''));

    expect(state.toJS()).toEqual({
      srv: {
        chan1: {
          users: [
            { mode: 'o', nick: 'nick2', renderName: '@nick2' },
            { mode: 'v', nick: 'nick1', renderName: '+nick1' }
          ]
        },
        chan2: {
          users: [
            { mode: '', nick: 'nick2', renderName: 'nick2' }
          ]
        }
      }
    });
  });

  it('handles SOCKET_CHANNELS', () => {
    const state = reducer(undefined, {
      type: actions.socket.CHANNELS,
      data: [
        { server: 'srv', name: 'chan1', topic: 'the topic' },
        { server: 'srv', name: 'chan2' },
        { server: 'srv2', name: 'chan1' }
      ]
    });

    expect(state.toJS()).toEqual({
      srv: {
        chan1: { topic: 'the topic', users: [] },
        chan2: { users: [] }
      },
      srv2: {
        chan1: { users: [] }
      }
    });
  });

  it('handles SOCKET_SERVERS', () => {
    const state = reducer(undefined, {
      type: actions.socket.SERVERS,
      data: [
        { host: '127.0.0.1' },
        { host: 'thehost' }
      ]
    });

    expect(state.toJS()).toEqual({
      '127.0.0.1': {},
      thehost: {}
    });
  });

  it('optimistically adds the server on CONNECT', () => {
    const state = reducer(undefined, connect('127.0.0.1:1337', 'nick', {}));

    expect(state.toJS()).toEqual({
      '127.0.0.1': {}
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
});

function socket_join(server, channel, user) {
  return {
    type: actions.socket.JOIN,
    server, user,
    channels: [channel]
  };
}

function socket_mode(server, channel, user, add, remove) {
  return {
    type: actions.socket.MODE,
    server, channel, user, add, remove
  };
}
