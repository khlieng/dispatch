import reducer, { compareUsers, getSortedChannels } from '../channels';
import { connect } from '../servers';
import * as actions from '../actions';

describe('channel reducer', () => {
  it('removes channels on PART', () => {
    let state = {
      srv1: {
        chan1: {},
        chan2: {},
        chan3: {}
      },
      srv2: {
        chan1: {}
      }
    };

    state = reducer(state, {
      type: actions.PART,
      server: 'srv1',
      channels: ['chan1', 'chan3']
    });

    expect(state).toEqual({
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

    expect(state).toEqual({
      srv: {
        chan1: {
          joined: true,
          users: [{ mode: '', nick: 'nick1', renderName: 'nick1' }]
        },
        chan2: {
          joined: true,
          users: [{ mode: '', nick: 'nick2', renderName: 'nick2' }]
        }
      }
    });
  });

  it('handles SOCKET_JOIN', () => {
    const state = reducer(undefined, socket_join('srv', 'chan1', 'nick1'));

    expect(state).toEqual({
      srv: {
        chan1: {
          joined: true,
          users: [{ mode: '', nick: 'nick1', renderName: 'nick1' }]
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

    expect(state).toEqual({
      srv: {
        chan1: {
          joined: true,
          users: [{ mode: '', nick: 'nick1', renderName: 'nick1' }]
        },
        chan2: {
          joined: true,
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
      oldNick: 'nick1',
      newNick: 'nick3'
    });

    expect(state).toEqual({
      srv: {
        chan1: {
          joined: true,
          users: [
            { mode: '', nick: 'nick3', renderName: 'nick3' },
            { mode: '', nick: 'nick2', renderName: 'nick2' }
          ]
        },
        chan2: {
          joined: true,
          users: [{ mode: '', nick: 'nick2', renderName: 'nick2' }]
        }
      }
    });
  });

  it('handles SOCKET_USERS', () => {
    let state = reducer(undefined, socket_join('srv', 'chan1', 'nick1'));
    state = reducer(state, {
      type: actions.socket.USERS,
      server: 'srv',
      channel: 'chan1',
      users: ['user3', 'user2', '@user4', 'user1', '+user5']
    });

    expect(state).toEqual({
      srv: {
        chan1: {
          joined: true,
          users: [
            { mode: '', nick: 'user3', renderName: 'user3' },
            { mode: '', nick: 'user2', renderName: 'user2' },
            { mode: 'o', nick: 'user4', renderName: '@user4' },
            { mode: '', nick: 'user1', renderName: 'user1' },
            { mode: 'v', nick: 'user5', renderName: '+user5' }
          ]
        }
      }
    });
  });

  it('handles SOCKET_TOPIC', () => {
    let state = reducer(undefined, socket_join('srv', 'chan1', 'nick1'));
    state = reducer(state, {
      type: actions.socket.TOPIC,
      server: 'srv',
      channel: 'chan1',
      topic: 'the topic'
    });

    expect(state).toMatchObject({
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

    expect(state).toMatchObject({
      srv: {
        chan1: {
          users: [
            { mode: 'o', nick: 'nick1', renderName: '@nick1' },
            { mode: '', nick: 'nick2', renderName: 'nick2' }
          ]
        },
        chan2: {
          users: [{ mode: '', nick: 'nick2', renderName: 'nick2' }]
        }
      }
    });

    state = reducer(state, socket_mode('srv', 'chan1', 'nick1', 'v', 'o'));
    state = reducer(state, socket_mode('srv', 'chan1', 'nick2', 'o', ''));
    state = reducer(state, socket_mode('srv', 'chan2', 'not_there', 'x', ''));

    expect(state).toMatchObject({
      srv: {
        chan1: {
          users: [
            { mode: 'v', nick: 'nick1', renderName: '+nick1' },
            { mode: 'o', nick: 'nick2', renderName: '@nick2' }
          ]
        },
        chan2: {
          users: [{ mode: '', nick: 'nick2', renderName: 'nick2' }]
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

    expect(state).toEqual({
      srv: {
        chan1: { joined: true, topic: 'the topic', users: [] },
        chan2: { joined: true, users: [] }
      },
      srv2: {
        chan1: { joined: true, users: [] }
      }
    });
  });

  it('handles SOCKET_SERVERS', () => {
    const state = reducer(undefined, {
      type: actions.socket.SERVERS,
      data: [{ host: '127.0.0.1' }, { host: 'thehost' }]
    });

    expect(state).toEqual({
      '127.0.0.1': {},
      thehost: {}
    });
  });

  it('optimistically adds the server on CONNECT', () => {
    const state = reducer(
      undefined,
      connect({ host: '127.0.0.1', nick: 'nick' })
    );

    expect(state).toEqual({
      '127.0.0.1': {}
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
});

function socket_join(server, channel, user) {
  return {
    type: actions.socket.JOIN,
    server,
    user,
    channels: [channel]
  };
}

function socket_mode(server, channel, user, add, remove) {
  return {
    type: actions.socket.MODE,
    server,
    channel,
    user,
    add,
    remove
  };
}

describe('compareUsers()', () => {
  it('compares users correctly', () => {
    expect(
      [
        { renderName: 'user5' },
        { renderName: '@user2' },
        { renderName: 'user3' },
        { renderName: 'user2' },
        { renderName: '+user1' },
        { renderName: '~bob' },
        { renderName: '%apples' },
        { renderName: '&cake' }
      ].sort(compareUsers)
    ).toEqual([
      { renderName: '~bob' },
      { renderName: '&cake' },
      { renderName: '@user2' },
      { renderName: '%apples' },
      { renderName: '+user1' },
      { renderName: 'user2' },
      { renderName: 'user3' },
      { renderName: 'user5' }
    ]);
  });
});

describe('getSortedChannels', () => {
  it('sorts servers and channels', () => {
    expect(
      getSortedChannels({
        channels: {
          'bob.com': {},
          '127.0.0.1': {
            '#chan1': {
              users: [],
              topic: 'cake'
            },
            '#pie': {},
            '##apples': {}
          }
        }
      })
    ).toEqual([
      {
        address: '127.0.0.1',
        channels: ['##apples', '#chan1', '#pie']
      },
      {
        address: 'bob.com',
        channels: []
      }
    ]);
  });
});
