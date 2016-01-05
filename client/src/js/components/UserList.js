import React, { Component } from 'react';
import Infinite from 'react-infinite';
import pure from 'pure-render-decorator';
import UserListItem from './UserListItem';

@pure
export default class UserList extends Component {
  state = {
    height: window.innerHeight - 100
  }

  componentDidMount() {
    window.addEventListener('resize', this.handleResize);
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.handleResize);
  }

  handleResize = () => {
    this.setState({ height: window.innerHeight - 100 });
  }

  render() {
    const { tab, openPrivateChat, select } = this.props;
    const users = [];
    const style = {};

    if (!tab.channel) {
      style.display = 'none';
    } else {
      this.props.users.forEach(user => users.push(
        <UserListItem
          key={user.nick}
          user={user}
          tab={tab}
          openPrivateChat={openPrivateChat}
          select={select}
        />
      ));
    }

    return (
      <div className="userlist" style={style}>
        <Infinite containerHeight={this.state.height} elementHeight={24}>
          {users}
        </Infinite>
      </div>
    );
  }
}
