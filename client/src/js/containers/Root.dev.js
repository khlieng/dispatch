import React, { Component } from 'react';
import { Router } from 'react-router';
import { Provider } from 'react-redux';
import pure from 'pure-render-decorator';
import DevTools from './DevTools';

@pure
export default class Root extends Component {
  render() {
    const { store, routes, history } = this.props;
    return (
      <Provider store={store}>
        <div>
          <Router routes={routes} history={history} />
          <DevTools />
        </div>
      </Provider>
    );
  }
}
