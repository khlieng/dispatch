import reducer, { compareUsers, getSortedChannels } from '../channels';
import { connect } from '../networks';
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
      network: 'srv1',
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
      network: 'srv',
      channel: 'chan1',
      user: 'nick2'
    });

    expect(state).toEqual({
      srv: {
        chan1: {
          name: 'chan1',
          joined: true,
          users: [{ mode: '', nick: 'nick1', renderName: 'nick1' }]
        },
        chan2: {
          name: 'chan2',
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
          name: 'chan1',
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
      network: 'srv',
      user: 'nick2'
    });

    expect(state).toEqual({
      srv: {
        chan1: {
          name: 'chan1',
          joined: true,
          users: [{ mode: '', nick: 'nick1', renderName: 'nick1' }]
        },
        chan2: {
          name: 'chan2',
          joined: true,
          users: []
        }
      }
    });
  });

  it('handles KICKED', () => {
    let state = reducer(
      undefined,
      connect({
        host: 'srv',
        nick: 'nick2'
      })
    );
    state = reducer(state, socket_join('srv', 'chan1', 'nick1'));
    state = reducer(state, socket_join('srv', 'chan1', 'nick2'));
    state = reducer(state, socket_join('srv', 'chan2', 'nick2'));

    state = reducer(state, {
      type: actions.KICKED,
      network: 'srv',
      channel: 'chan2',
      user: 'nick2',
      self: true
    });

    expect(state).toEqual({
      srv: {
        chan1: {
          name: 'chan1',
          joined: true,
          users: [
            { mode: '', nick: 'nick1', renderName: 'nick1' },
            { mode: '', nick: 'nick2', renderName: 'nick2' }
          ]
        },
        chan2: {
          name: 'chan2',
          joined: false,
          users: []
        }
      }
    });

    state = reducer(state, {
      type: actions.KICKED,
      network: 'srv',
      channel: 'chan1',
      user: 'nick1'
    });

    expect(state).toEqual({
      srv: {
        chan1: {
          name: 'chan1',
          joined: true,
          users: [{ mode: '', nick: 'nick2', renderName: 'nick2' }]
        },
        chan2: {
          name: 'chan2',
          joined: false,
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
      network: 'srv',
      oldNick: 'nick1',
      newNick: 'nick3'
    });

    expect(state).toEqual({
      srv: {
        chan1: {
          name: 'chan1',
          joined: true,
          users: [
            { mode: '', nick: 'nick3', renderName: 'nick3' },
            { mode: '', nick: 'nick2', renderName: 'nick2' }
          ]
        },
        chan2: {
          name: 'chan2',
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
      network: 'srv',
      channel: 'chan1',
      users: ['user3', 'user2', '@user4', 'user1', '+user5']
    });

    expect(state).toEqual({
      srv: {
        chan1: {
          name: 'chan1',
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
      network: 'srv',
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

  it('handles channels from INIT', () => {
    const state = reducer(undefined, {
      type: actions.INIT,
      channels: [
        { network: 'srv', name: 'chan1', topic: 'the topic' },
        { network: 'srv', name: 'chan2', joined: true },
        { network: 'srv2', name: 'chan1' }
      ]
    });

    expect(state).toEqual({
      srv: {
        chan1: { name: 'chan1', topic: 'the topic', users: [] },
        chan2: { name: 'chan2', joined: true, users: [] }
      },
      srv2: {
        chan1: { name: 'chan1', users: [] }
      }
    });
  });

  it('handles networks from INIT', () => {
    const state = reducer(undefined, {
      type: actions.INIT,
      networks: [{ host: '127.0.0.1' }, { host: 'thehost' }]
    });

    expect(state).toEqual({
      '127.0.0.1': {},
      thehost: {}
    });
  });

  it('optimistically adds the network on CONNECT', () => {
    const state = reducer(
      undefined,
      connect({ host: '127.0.0.1', nick: 'nick' })
    );

    expect(state).toEqual({
      '127.0.0.1': {}
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
});

function socket_join(network, channel, user) {
  return {
    type: actions.socket.JOIN,
    network,
    user,
    channels: [channel]
  };
}

function socket_mode(network, channel, user, add, remove) {
  return {
    type: actions.socket.MODE,
    network,
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
  it('sorts networks and channels', () => {
    expect(
      getSortedChannels({
        channels: {
          'bob.com': {},
          '127.0.0.1': {
            '#chan1': {
              name: '#chan1',
              users: [],
              topic: 'cake'
            },
            '#pie': {
              name: '#pie'
            },
            '##apples': {
              name: '##apples'
            }
          }
        }
      })
    ).toEqual([
      {
        address: '127.0.0.1',
        channels: [
          {
            name: '##apples'
          },
          {
            name: '#chan1',
            users: [],
            topic: 'cake'
          },
          {
            name: '#pie'
          }
        ]
      },
      {
        address: 'bob.com',
        channels: []
      }
    ]);
  });
});
