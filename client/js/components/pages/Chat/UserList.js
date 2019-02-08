import React, { PureComponent, createRef } from 'react';
import { VariableSizeList as List } from 'react-window';
import AutoSizer from 'react-virtualized-auto-sizer';
import classnames from 'classnames';
import UserListItem from './UserListItem';

export default class UserList extends PureComponent {
  list = createRef();

  getSnapshotBeforeUpdate(prevProps) {
    if (this.list.current) {
      const { users } = this.props;

      if (prevProps.users.length !== users.length) {
        this.list.current.resetAfterIndex(
          Math.min(prevProps.users.length, users.length) + 1
        );
      } else {
        this.list.current.forceUpdate();
      }
    }

    return null;
  }

  getItemHeight = index => {
    const { users } = this.props;

    if (index === 0) {
      return 12;
    } if (index === users.length + 1) {
      return 10;
    }
    return 24;
  };

  getItemKey = index => {
    const { users } = this.props;

    if (index === 0) {
      return 'top';
    } if (index === users.length + 1) {
      return 'bottom';
    }
    return index;
  };

  renderUser = ({ index, style }) => {
    const { users, coloredNicks, onNickClick } = this.props;

    if (index === 0 || index === users.length + 1) {
      return null;
    }

    return (
      <UserListItem
        user={users[index - 1]}
        coloredNick={coloredNicks}
        style={style}
        onClick={onNickClick}
      />
    );
  };

  render() {
    const { users, showUserList } = this.props;

    const className = classnames('userlist', {
      'off-canvas': showUserList
    });

    return (
      <div className={className}>
        <AutoSizer disableWidth>
          {({ height }) => (
            <List
              ref={this.list}
              width={200}
              height={height}
              itemCount={users.length + 2}
              itemKey={this.getItemKey}
              itemSize={this.getItemHeight}
              estimatedItemSize={24}
              overscanCount={5}
            >
              {this.renderUser}
            </List>
          )}
        </AutoSizer>
      </div>
    );
  }
}
