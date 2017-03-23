import React, { PureComponent } from 'react';
import { List } from 'react-virtualized/dist/commonjs/List';
import { AutoSizer } from 'react-virtualized/dist/commonjs/AutoSizer';
import UserListItem from './UserListItem';

const listStyle = { padding: '10px 0', boxSizing: 'content-box' };

export default class UserList extends PureComponent {
  componentWillUpdate(nextProps) {
    if (nextProps.users.size === this.props.users.size) {
      this.list.forceUpdate();
    }
  }

  listRef = el => { this.list = el; };

  renderUser = ({ index, style }) => {
    const { users, tab, openPrivateChat, select } = this.props;
    const user = users.get(index);

    return (
      <UserListItem
        key={user.nick}
        user={user}
        tab={tab}
        openPrivateChat={openPrivateChat}
        select={select}
        style={style}
      />
    );
  };

  render() {
    const { tab, showUserList } = this.props;
    const className = showUserList ? 'userlist off-canvas' : 'userlist';
    const style = {};

    if (!tab.channel) {
      style.display = 'none';
    }

    return (
      <div className={className} style={style}>
        <AutoSizer disableWidth>
          {({ height }) => (
            <List
              ref={this.listRef}
              width={200}
              height={height - 20}
              rowCount={this.props.users.size}
              rowHeight={24}
              rowRenderer={this.renderUser}
              style={listStyle}
            />
          )}
        </AutoSizer>
      </div>
    );
  }
}
