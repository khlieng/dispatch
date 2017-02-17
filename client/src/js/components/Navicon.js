import React, { PureComponent } from 'react';
import { connect } from 'react-redux';
import { toggleMenu } from '../actions/ui';

class Navicon extends PureComponent {
  handleClick = () => this.props.dispatch(toggleMenu());

  render() {
    return (
      <i className="icon-menu navicon" onClick={this.handleClick} />
    );
  }
}

export default connect()(Navicon);
