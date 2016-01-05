import React, { Component } from 'react';
import pure from 'pure-render-decorator';

@pure
export default class TabListItem extends Component {
  handleClick = () => {
    const { server, target, onClick } = this.props;
    onClick(server, target);
  }

  render() {
    const { target, content, selected } = this.props;
    const classes = [];

    if (!target) {
      classes.push('tab-server');
    }

    if (selected) {
      classes.push('selected');
    }

    return (
      <p className={classes.join(' ')} onClick={this.handleClick}>{content}</p>
    );
  }
}
