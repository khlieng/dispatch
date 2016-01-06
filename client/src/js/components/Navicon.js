import React, { Component } from 'react';
import { connect } from 'react-redux';
import pure from 'pure-render-decorator';
import { toggleMenu } from '../actions/tab';

@pure
class Navicon extends Component {
  handleClick = () => this.props.dispatch(toggleMenu());

  render() {
    return (
      <i className="icon-menu navicon" onClick={this.handleClick}></i>
    );
  }
}

export default connect()(Navicon);
