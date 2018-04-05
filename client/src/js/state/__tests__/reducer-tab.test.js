import reducer, { setSelectedTab } from '../tab';
import * as actions from '../actions';
import { locationChanged } from 'utils/router';

describe('tab reducer', () => {
  it('selects the tab and adds it to history', () => {
    let state = reducer(undefined, setSelectedTab('srv', '#chan'));

    expect(state.toJS()).toEqual({
      selected: { server: 'srv', name: '#chan' },
      history: [
        { server: 'srv', name: '#chan' }
      ]
    });

    state = reducer(state, setSelectedTab('srv', 'user1'));

    expect(state.toJS()).toEqual({
      selected: { server: 'srv', name: 'user1' },
      history: [
        { server: 'srv', name: '#chan' },
        { server: 'srv', name: 'user1' }
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
      server: 'srv',
      channels: ['#chan']
    });

    expect(state.toJS()).toEqual({
      selected: { server: 'srv', name: '#chan3' },
      history: [
        { server: 'srv1', name: 'bob' },
        { server: 'srv', name: '#chan3' }
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
      server: 'srv1',
      nick: 'bob'
    });

    expect(state.toJS()).toEqual({
      selected: { server: 'srv', name: '#chan3' },
      history: [
        { server: 'srv', name: '#chan' },
        { server: 'srv', name: '#chan' },
        { server: 'srv', name: '#chan3' }
      ]
    });
  });

  it('removes all tabs related to server from history on DISCONNECT', () => {
    let state = reducer(undefined, setSelectedTab('srv', '#chan'));
    state = reducer(state, setSelectedTab('srv1', 'bob'));
    state = reducer(state, setSelectedTab('srv', '#chan'));
    state = reducer(state, setSelectedTab('srv', '#chan3'));

    state = reducer(state, {
      type: actions.DISCONNECT,
      server: 'srv',
    });

    expect(state.toJS()).toEqual({
      selected: { server: 'srv', name: '#chan3' },
      history: [
        { server: 'srv1', name: 'bob' },
      ]
    });
  });

  it('clears the tab when navigating to a non-tab page', () => {
    let state = reducer(undefined, setSelectedTab('srv', '#chan'));

    state = reducer(state, locationChanged('settings'));

    expect(state.toJS()).toEqual({
      selected: { server: null, name: null },
      history: [
        { server: 'srv', name: '#chan' }
      ]
    });
  });

  it('selects the tab and adds it to history when navigating to a tab', () => {
    const state = reducer(undefined,
      locationChanged('chat', {
        server: 'srv',
        name: '#chan'
      })
    );

    expect(state.toJS()).toEqual({
      selected: { server: 'srv', name: '#chan' },
      history: [
        { server: 'srv', name: '#chan' }
      ]
    });
  });
});
