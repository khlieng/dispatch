import React, { Component } from 'react';
import pure from 'pure-render-decorator';

@pure
export default class UserListItem extends Component {
  handleClick = () => {
    const { tab, user, openPrivateChat, select } = this.props;

    openPrivateChat(tab.server, user.nick);
    select(tab.server, user.nick, true);
  };

  render() {
    return <p onClick={this.handleClick}>{this.props.user.renderName}</p>;
  }
}
