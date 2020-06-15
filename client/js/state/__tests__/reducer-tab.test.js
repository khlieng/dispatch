import reducer, { setSelectedTab } from '../tab';
import * as actions from '../actions';
import { locationChanged } from 'utils/router';

describe('tab reducer', () => {
  it('selects the tab and adds it to history', () => {
    let state = reducer(undefined, setSelectedTab('srv', '#chan'));

    expect(state).toEqual({
      selected: { network: 'srv', name: '#chan' },
      history: [{ network: 'srv', name: '#chan' }]
    });

    state = reducer(state, setSelectedTab('srv', 'user1'));

    expect(state).toEqual({
      selected: { network: 'srv', name: 'user1' },
      history: [
        { network: 'srv', name: '#chan' },
        { network: 'srv', name: 'user1' }
      ]
    });
  });

  it('removes the tab from history on PART', () => {
    let state = reducer(undefined, setSelectedTab('srv', '#chan'));
    state = reducer(state, setSelectedTab('srv1', 'bob'));
    state = reducer(state, setSelectedTab('srv', '#chan'));
    state = reducer(state, setSelectedTab('srv', '#chan3'));

    state = reducer(state, {
      type: actions.PART,
      network: 'srv',
      channels: ['#chan']
    });

    expect(state).toEqual({
      selected: { network: 'srv', name: '#chan3' },
      history: [
        { network: 'srv1', name: 'bob' },
        { network: 'srv', name: '#chan3' }
      ]
    });
  });

  it('removes the tab from history on CLOSE_PRIVATE_CHAT', () => {
    let state = reducer(undefined, setSelectedTab('srv', '#chan'));
    state = reducer(state, setSelectedTab('srv1', 'bob'));
    state = reducer(state, setSelectedTab('srv', '#chan'));
    state = reducer(state, setSelectedTab('srv', '#chan3'));

    state = reducer(state, {
      type: actions.CLOSE_PRIVATE_CHAT,
      network: 'srv1',
      nick: 'bob'
    });

    expect(state).toEqual({
      selected: { network: 'srv', name: '#chan3' },
      history: [
        { network: 'srv', name: '#chan' },
        { network: 'srv', name: '#chan' },
        { network: 'srv', name: '#chan3' }
      ]
    });
  });

  it('removes all tabs related to network from history on DISCONNECT', () => {
    let state = reducer(undefined, setSelectedTab('srv', '#chan'));
    state = reducer(state, setSelectedTab('srv1', 'bob'));
    state = reducer(state, setSelectedTab('srv', '#chan'));
    state = reducer(state, setSelectedTab('srv', '#chan3'));

    state = reducer(state, {
      type: actions.DISCONNECT,
      network: 'srv'
    });

    expect(state).toEqual({
      selected: { network: 'srv', name: '#chan3' },
      history: [{ network: 'srv1', name: 'bob' }]
    });
  });

  it('clears the tab when navigating to a non-tab page', () => {
    let state = reducer(undefined, setSelectedTab('srv', '#chan'));

    state = reducer(state, locationChanged('settings', {}, {}));

    expect(state).toEqual({
      selected: {},
      history: [{ network: 'srv', name: '#chan' }]
    });
  });

  it('selects the tab and adds it to history when navigating to a tab', () => {
    const state = reducer(
      undefined,
      locationChanged(
        'chat',
        {
          network: 'srv',
          name: '#chan'
        },
        {}
      )
    );

    expect(state).toEqual({
      selected: { network: 'srv', name: '#chan' },
      history: [{ network: 'srv', name: '#chan' }]
    });
  });
});
