import React, { PureComponent } from 'react';
import { List } from 'react-virtualized/dist/commonjs/List';
import { AutoSizer } from 'react-virtualized/dist/commonjs/AutoSizer';
import UserListItem from './UserListItem';

export default class UserList extends PureComponent {
  componentWillUpdate(nextProps) {
    if (nextProps.users.size === this.props.users.size) {
      this.list.forceUpdateGrid();
    }
  }

  listRef = el => { this.list = el; };

  renderUser = ({ index, style, key }) => {
    const { users, onNickClick } = this.props;

    return (
      <UserListItem
        key={key}
        user={users.get(index)}
        style={style}
        onClick={onNickClick}
      />
    );
  };

  render() {
    const { users, showUserList } = this.props;
    const className = showUserList ? 'userlist off-canvas' : 'userlist';

    return (
      <div className={className}>
        <AutoSizer disableWidth>
          {({ height }) => (
            <List
              ref={this.listRef}
              width={200}
              height={height - 20}
              rowCount={users.size}
              rowHeight={24}
              rowRenderer={this.renderUser}
              className="rvlist-users"
            />
          )}
        </AutoSizer>
      </div>
    );
  }
}
