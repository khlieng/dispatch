import reducer, { broadcast, getMessageTab } from '../messages';
import * as actions from '../actions';
import appReducer from '../app';

describe('message reducer', () => {
  it('adds the message on ADD_MESSAGE', () => {
    const state = reducer(undefined, {
      type: actions.ADD_MESSAGE,
      server: 'srv',
      tab: '#chan1',
      message: {
        from: 'foo',
        content: 'msg'
      }
    });

    expect(state).toMatchObject({
      srv: {
        '#chan1': [
          {
            from: 'foo',
            content: 'msg'
          }
        ]
      }
    });
  });

  it('adds all the messages on ADD_MESSAGES', () => {
    const state = reducer(undefined, {
      type: actions.ADD_MESSAGES,
      server: 'srv',
      tab: '#chan1',
      messages: [
        {
          from: 'foo',
          content: 'msg'
        },
        {
          from: 'bar',
          content: 'msg'
        },
        {
          tab: '#chan2',
          from: 'foo',
          content: 'msg'
        }
      ]
    });

    expect(state).toMatchObject({
      srv: {
        '#chan1': [
          {
            from: 'foo',
            content: 'msg'
          },
          {
            from: 'bar',
            content: 'msg'
          }
        ],
        '#chan2': [
          {
            from: 'foo',
            content: 'msg'
          }
        ]
      }
    });
  });

  it('handles prepending of messages on ADD_MESSAGES', () => {
    let state = {
      srv: {
        '#chan1': [{ id: 0 }]
      }
    };

    state = reducer(state, {
      type: actions.ADD_MESSAGES,
      server: 'srv',
      tab: '#chan1',
      prepend: true,
      messages: [
        { id: 1, date: new Date() },
        { id: 2, date: new Date() }
      ]
    });

    expect(state).toMatchObject({
      srv: {
        '#chan1': [{ id: 1 }, { id: 2 }, { id: 0 }]
      }
    });
  });

  it('adds date markers when prepending messages', () => {
    let state = {
      srv: {
        '#chan1': [{ id: 0, date: new Date(1999, 0, 1) }]
      }
    };

    state = reducer(state, {
      type: actions.ADD_MESSAGES,
      server: 'srv',
      tab: '#chan1',
      prepend: true,
      messages: [
        { id: 1, date: new Date(1990, 0, 2) },
        { id: 2, date: new Date(1990, 0, 3) }
      ]
    });

    expect(state).toMatchObject({
      srv: {
        '#chan1': [
          { id: 1 },
          { type: 'date' },
          { id: 2 },
          { type: 'date' },
          { id: 0 }
        ]
      }
    });
  });

  it('adds a date marker when adding a message', () => {
    let state = {
      srv: {
        '#chan1': [{ id: 0, date: new Date(1999, 0, 1) }]
      }
    };

    state = reducer(state, {
      type: actions.ADD_MESSAGE,
      server: 'srv',
      tab: '#chan1',
      message: { id: 1, date: new Date(1990, 0, 2) }
    });

    expect(state).toMatchObject({
      srv: {
        '#chan1': [{ id: 0 }, { type: 'date' }, { id: 1 }]
      }
    });
  });

  it('adds date markers when adding messages', () => {
    let state = {
      srv: {
        '#chan1': [{ id: 0, date: new Date(1999, 0, 1) }]
      }
    };

    state = reducer(state, {
      type: actions.ADD_MESSAGES,
      server: 'srv',
      tab: '#chan1',
      messages: [
        { id: 1, date: new Date(1990, 0, 2) },
        { id: 2, date: new Date(1990, 0, 3) },
        { id: 3, date: new Date(1990, 0, 3) }
      ]
    });

    expect(state).toMatchObject({
      srv: {
        '#chan1': [
          { id: 0 },
          { type: 'date' },
          { id: 1 },
          { type: 'date' },
          { id: 2 },
          { id: 3 }
        ]
      }
    });
  });

  it('adds messages to the correct tabs when broadcasting', () => {
    let state = {
      app: appReducer(undefined, { type: '' })
    };

    const thunk = broadcast('test', 'srv', ['#chan1', '#chan3']);
    thunk(
      action => {
        state.messages = reducer(undefined, action);
      },
      () => state
    );

    const messages = state.messages;

    expect(messages.srv).not.toHaveProperty('srv');
    expect(messages.srv['#chan1']).toHaveLength(1);
    expect(messages.srv['#chan1'][0].content).toBe('test');
    expect(messages.srv['#chan3']).toHaveLength(1);
    expect(messages.srv['#chan3'][0].content).toBe('test');
  });

  it('deletes all messages related to server when disconnecting', () => {
    let state = {
      srv: {
        '#chan1': [{ content: 'msg1' }, { content: 'msg2' }],
        '#chan2': [{ content: 'msg' }]
      },
      srv2: {
        '#chan1': [{ content: 'msg' }]
      }
    };

    state = reducer(state, {
      type: actions.DISCONNECT,
      server: 'srv'
    });

    expect(state).toEqual({
      srv2: {
        '#chan1': [{ content: 'msg' }]
      }
    });
  });

  it('deletes all messages related to channel when parting', () => {
    let state = {
      srv: {
        '#chan1': [{ content: 'msg1' }, { content: 'msg2' }],
        '#chan2': [{ content: 'msg' }]
      },
      srv2: {
        '#chan1': [{ content: 'msg' }]
      }
    };

    state = reducer(state, {
      type: actions.PART,
      server: 'srv',
      channels: ['#chan1']
    });

    expect(state).toEqual({
      srv: {
        '#chan2': [{ content: 'msg' }]
      },
      srv2: {
        '#chan1': [{ content: 'msg' }]
      }
    });
  });

  it('deletes direct messages when closing a direct message tab', () => {
    let state = {
      srv: {
        bob: [{ content: 'msg1' }, { content: 'msg2' }],
        '#chan2': [{ content: 'msg' }]
      },
      srv2: {
        '#chan1': [{ content: 'msg' }]
      }
    };

    state = reducer(state, {
      type: actions.CLOSE_PRIVATE_CHAT,
      server: 'srv',
      nick: 'bob'
    });

    expect(state).toEqual({
      srv: {
        '#chan2': [{ content: 'msg' }]
      },
      srv2: {
        '#chan1': [{ content: 'msg' }]
      }
    });
  });
});

describe('getMessageTab()', () => {
  it('returns the correct tab', () => {
    const srv = 'chat.freenode.net';
    [
      ['#cake', '#cake'],
      ['#apple.pie', '#apple.pie'],
      ['bob', 'bob'],
      [undefined, srv],
      [null, srv],
      ['*', srv],
      [srv, srv],
      ['beans.freenode.net', srv]
    ].forEach(([target, expected]) =>
      expect(getMessageTab(srv, target)).toBe(expected)
    );
  });
});
