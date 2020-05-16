import React, { Suspense, lazy, useState, useEffect } from 'react';
import Route from 'containers/Route';
import AppInfo from 'components/AppInfo';
import TabList from 'components/TabList';
import cn from 'classnames';

const Modals = lazy(() => import('components/modals'));
const Chat = lazy(() => import('containers/Chat'));
const Connect = lazy(() =>
  import(/* webpackChunkName: "connect" */ 'containers/Connect')
);
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
  hideMenu,
  openModal,
  newVersionAvailable,
  hasOpenModals
}) => {
  const [renderModals, setRenderModals] = useState(false);
  if (!renderModals && hasOpenModals) {
    setRenderModals(true);
  }

  const [starting, setStarting] = useState(true);
  useEffect(() => {
    setTimeout(() => setStarting(false), 1000);
  }, []);

  const mainClass = cn('main-container', {
    'off-canvas': showTabList
  });

  const handleClick = () => {
    if (showTabList) {
      hideMenu();
    }
  };

  return (
    <div className="wrap" onClick={handleClick}>
      {!starting && !connected && (
        <AppInfo type="error">
          Connection lost, attempting to reconnect...
        </AppInfo>
      )}
      {newVersionAvailable && (
        <AppInfo dismissible>
          A new version of dispatch just got installed, reload to start using
          it!
        </AppInfo>
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
          openModal={openModal}
        />
        <div className={mainClass}>
          <Suspense fallback={<div className="suspense-fallback">...</div>}>
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
          <Suspense
            fallback={<div className="suspense-modal-fallback">...</div>}
          >
            {renderModals && <Modals />}
          </Suspense>
        </div>
      </div>
    </div>
  );
};

export default App;
