import React from 'react';
import Route from '../containers/Route';
import Chat from '../containers/Chat';
import Connect from '../containers/Connect';
import Settings from '../containers/Settings';
import TabList from '../components/TabList';

const App = props => {
  const { onClick, ...tabListProps } = props;
  const mainClass = props.showTabList ? 'main-container off-canvas' : 'main-container';

  return (
    <div onClick={onClick}>
      <TabList {...tabListProps} />
      <div className={mainClass}>
        <Route name="chat"><Chat /></Route>
        <Route name="connect"><Connect /></Route>
        <Route name="settings"><Settings /></Route>
      </div>
    </div>
  );
};

export default App;
