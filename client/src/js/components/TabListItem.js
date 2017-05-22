import React, { PureComponent } from 'react';

export default class TabListItem extends PureComponent {
  handleClick = () => {
    const { server, target, onClick } = this.props;
    onClick(server, target);
  };

  render() {
    const { target, content, selected, connected } = this.props;
    const classes = [];
    const style = {};

    if (!target) {
      classes.push('tab-server');

      if (connected) {
        style.color = '#6BB758';
      } else {
        style.color = '#F6546A';
      }
    }

    if (selected) {
      classes.push('selected');
    }

    return (
      <p className={classes.join(' ')} style={style} onClick={this.handleClick}>
        <span className="tab-content">{content}</span>
      </p>
    );
  }
}
