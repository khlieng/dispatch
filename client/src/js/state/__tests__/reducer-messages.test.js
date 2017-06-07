import { Map, fromJS } from 'immutable';
import reducer, { broadcast } from '../messages';
import * as actions from '../actions';
import appReducer from '../app';

describe('reducers/messages', () => {
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

    expect(state.toJS()).toMatchObject({
      srv: {
        '#chan1': [{
          from: 'foo',
          content: 'msg'
        }]
      }
    });
  });

  it('adds all the messsages on ADD_MESSAGES', () => {
    const state = reducer(undefined, {
      type: actions.ADD_MESSAGES,
      server: 'srv',
      tab: '#chan1',
      messages: [
        {
          from: 'foo',
          content: 'msg'
        }, {
          from: 'bar',
          content: 'msg'
        }, {
          tab: '#chan2',
          from: 'foo',
          content: 'msg'
        }
      ]
    });

    expect(state.toJS()).toMatchObject({
      srv: {
        '#chan1': [
          {
            from: 'foo',
            content: 'msg'
          }, {
            from: 'bar',
            content: 'msg'
          }
        ],
        '#chan2': [{
          from: 'foo',
          content: 'msg'
        }]
      }
    });
  });

  it('handles prepending of messages on ADD_MESSAGES', () => {
    let state = fromJS({
      srv: {
        '#chan1': [{ id: 0 }]
      }
    });

    state = reducer(state, {
      type: actions.ADD_MESSAGES,
      server: 'srv',
      tab: '#chan1',
      prepend: true,
      messages: [{ id: 1 }, { id: 2 }]
    });

    expect(state.toJS()).toMatchObject({
      srv: {
        '#chan1': [{ id: 1 }, { id: 2 }, { id: 0 }]
      }
    });
  });

  it('adds messages to the correct tabs when broadcasting', () => {
    let state = {
      app: appReducer(undefined, { type: '' })
    };

    const thunk = broadcast('test', 'srv', ['#chan1', '#chan3']);
    thunk(
      action => { state.messages = reducer(undefined, action); },
      () => state
    );

    const messages = state.messages.toJS();

    expect(messages.srv).not.toHaveProperty('srv');
    expect(messages.srv['#chan1']).toHaveLength(1);
    expect(messages.srv['#chan1'][0].content).toBe('test');
    expect(messages.srv['#chan3']).toHaveLength(1);
    expect(messages.srv['#chan3'][0].content).toBe('test');
  });

  it('deletes all messages related to server when disconnecting', () => {
    let state = fromJS({
      srv: {
        '#chan1': [
          { content: 'msg1' },
          { content: 'msg2' }
        ],
        '#chan2': [
          { content: 'msg' }
        ]
      },
      srv2: {
        '#chan1': [
          { content: 'msg' }
        ]
      }
    });

    state = reducer(state, {
      type: actions.DISCONNECT,
      server: 'srv'
    });

    expect(state.toJS()).toEqual({
      srv2: {
        '#chan1': [
          { content: 'msg' }
        ]
      }
    });
  });

  it('deletes all messages related to channel when parting', () => {
    let state = fromJS({
      srv: {
        '#chan1': [
          { content: 'msg1' },
          { content: 'msg2' }
        ],
        '#chan2': [
          { content: 'msg' }
        ]
      },
      srv2: {
        '#chan1': [
          { content: 'msg' }
        ]
      }
    });

    state = reducer(state, {
      type: actions.PART,
      server: 'srv',
      channels: ['#chan1']
    });

    expect(state.toJS()).toEqual({
      srv: {
        '#chan2': [
          { content: 'msg' }
        ]
      },
      srv2: {
        '#chan1': [
          { content: 'msg' }
        ]
      }
    });
  });
});
