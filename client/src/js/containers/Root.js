import React, { Component } from 'react';
import { Router } from 'react-router';
import { Provider } from 'react-redux';

export default class Root extends Component {
  render() {
    const { store, routes, history } = this.props;
    return (
      <Provider store={store}>
        <Router routes={routes} history={history} />
      </Provider>
    );
  }
}
