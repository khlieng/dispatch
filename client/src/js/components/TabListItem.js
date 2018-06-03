import React, { PureComponent } from 'react';
import classnames from 'classnames';

export default class TabListItem extends PureComponent {
  handleClick = () => {
    const { server, target, onClick } = this.props;
    onClick(server, target);
  };

  render() {
    const { target, content, selected, connected } = this.props;

    const className = classnames({
      'tab-server': !target,
      success: !target && connected,
      error: !target && !connected,
      selected
    });

    return (
      <p className={className} onClick={this.handleClick}>
        <span className="tab-content">{content}</span>
      </p>
    );
  }
}
