import React, { Component } from 'react';
import pure from 'pure-render-decorator';

@pure
export default class TabListItem extends Component {
  render() {
    const classes = [];

    if (this.props.server) {
      classes.push('tab-server');
    }

    if (this.props.selected) {
      classes.push('selected');
    }

    return (
      <p className={classes.join(' ')} onClick={this.props.onClick}>{this.props.content}</p>
    );
  }
}
