import { Map } from 'immutable';
import reducer from '../reducers/messages';
import * as actions from '../actions';
import { broadcast } from '../actions/message';

describe('reducers/messages', () => {
  it('adds messages to the correct tabs when broadcasting', () => {
    let state = {
      environment: Map({
        charWidth: 0,
        wrapWidth: 0
      })
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
});
