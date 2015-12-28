import React, { Component } from 'react';
import { connect } from 'react-redux';
import pure from 'pure-render-decorator';
import { toggleMenu } from '../actions/tab';

@pure
class Navicon extends Component {
  render() {
    const { dispatch } = this.props;
    return (
      <i className="icon-menu navicon" onClick={() => dispatch(toggleMenu())}></i>
    );
  }
}

export default connect()(Navicon);
