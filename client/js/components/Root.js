import React from 'react';
import { Provider } from 'react-redux';
import { hot, setConfig } from 'react-hot-loader';
import App from 'containers/App';

setConfig({
  pureSFC: true
});

const Root = ({ store }) => (
  <Provider store={store}>
    <App />
  </Provider>
);

export default hot(module)(Root);
