import React from 'react';
import { Route, IndexRoute } from 'react-router';
import App from './containers/App';
import Connect from './containers/Connect';
import Chat from './containers/Chat';
import Settings from './containers/Settings';

export default function createRoutes() {
  return (
    <Route path="/" component={App}>
      <Route path="connect" component={Connect} />
      <Route path="settings" component={Settings} />
      <Route path="/:server" component={Chat} />
      <Route path="/:server/:channel" component={Chat} />
      <Route path="/:server/pm/:user" component={Chat} />
      <IndexRoute component={null} />
    </Route>
  );
}
