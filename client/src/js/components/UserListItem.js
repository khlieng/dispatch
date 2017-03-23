import React, { PureComponent } from 'react';

export default class UserListItem extends PureComponent {
  handleClick = () => {
    const { tab, user, openPrivateChat, select } = this.props;

    openPrivateChat(tab.server, user.nick);
    select(tab.server, user.nick, true);
  };

  render() {
    return (
      <p style={this.props.style} onClick={this.handleClick}>
        {this.props.user.renderName}
      </p>
    );
  }
}
