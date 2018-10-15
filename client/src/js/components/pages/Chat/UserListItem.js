import React, { PureComponent } from 'react';
import stringToRGB from 'utils/color';

export default class UserListItem extends PureComponent {
  handleClick = () => this.props.onClick(this.props.user.nick);

  render() {
    const { user, coloredNick } = this.props;
    let { style } = this.props;

    if (coloredNick) {
      style = {
        color: stringToRGB(user.nick),
        ...style
      };
    }

    return (
      <p style={style} onClick={this.handleClick}>
        {user.renderName}
      </p>
    );
  }
}
