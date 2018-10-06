import React, { PureComponent } from 'react';
import stringToHSL from 'utils/color';

export default class UserListItem extends PureComponent {
  handleClick = () => this.props.onClick(this.props.user.nick);

  render() {
    const { user } = this.props;
    const style = {
      color: stringToHSL(user.nick),
      ...this.props.style
    };

    return (
      <p style={style} onClick={this.handleClick}>
        {user.renderName}
      </p>
    );
  }
}
