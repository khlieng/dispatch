import React, { PureComponent } from 'react';

export default class UserListItem extends PureComponent {
  handleClick = () => this.props.onClick(this.props.user.nick);

  render() {
    return (
      <p style={this.props.style} onClick={this.handleClick}>
        {this.props.user.renderName}
      </p>
    );
  }
}
