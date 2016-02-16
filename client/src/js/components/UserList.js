import React, { Component } from 'react';
import { VirtualScroll } from 'react-virtualized';
import pure from 'pure-render-decorator';
import UserListItem from './UserListItem';

@pure
export default class UserList extends Component {
  state = {
    height: window.innerHeight - 100
  };

  componentDidMount() {
    window.addEventListener('resize', this.handleResize);
  }

  componentWillUpdate(nextProps) {
    if (nextProps.users.size === this.props.users.size) {
      this.refs.list.forceUpdate();
    }
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.handleResize);
  }

  getRowHeight = index => {
    if (index === 0 || index === this.props.users.size + 1) {
      return 10;
    }

    return 24;
  };

  handleResize = () => this.setState({ height: window.innerHeight - 100 });

  renderUser = index => {
    const { users } = this.props;

    if (index === 0 || index === users.size + 1) {
      return <span style={{ height: '10px' }}></span>;
    }

    const { tab, openPrivateChat, select } = this.props;
    const user = users.get(index - 1);
    return (
      <UserListItem
        key={user.nick}
        user={user}
        tab={tab}
        openPrivateChat={openPrivateChat}
        select={select}
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
        <VirtualScroll
          ref="list"
          height={this.state.height}
          rowsCount={this.props.users.size + 2}
          rowHeight={this.getRowHeight}
          rowRenderer={this.renderUser}
        />
      </div>
    );
  }
}
