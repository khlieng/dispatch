import React, { Suspense, lazy } from 'react';
import Route from 'containers/Route';
import TabList from 'components/TabList';
import classnames from 'classnames';

const Chat = lazy(() => import('containers/Chat'));
const Connect = lazy(() => import('containers/Connect'));
const Settings = lazy(() => import('containers/Settings'));

const App = ({
  connected,
  tab,
  channels,
  servers,
  privateChats,
  showTabList,
  select,
  push,
  hideMenu
}) => {
  const mainClass = classnames('main-container', {
    'off-canvas': showTabList
  });

  const handleClick = () => {
    if (showTabList) {
      hideMenu();
    }
  };

  return (
    <div className="wrap" onClick={handleClick}>
      {!connected && (
        <div className="app-info">
          Connection lost, attempting to reconnect...
        </div>
      )}
      <div className="app-container">
        <TabList
          tab={tab}
          channels={channels}
          servers={servers}
          privateChats={privateChats}
          showTabList={showTabList}
          select={select}
          push={push}
        />
        <div className={mainClass}>
          <Suspense
            maxDuration={1000}
            fallback={<div className="suspense-fallback">...</div>}
          >
            <Route name="chat">
              <Chat />
            </Route>
            <Route name="connect">
              <Connect />
            </Route>
            <Route name="settings">
              <Settings />
            </Route>
          </Suspense>
        </div>
      </div>
    </div>
  );
};

export default App;
