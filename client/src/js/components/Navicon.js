import React, { PureComponent } from 'react';
import { connect } from 'react-redux';
import { toggleMenu } from '../state/ui';

class Navicon extends PureComponent {
  render() {
    return (
      <i className="icon-menu navicon" onClick={this.props.toggleMenu} />
    );
  }
}

export default connect(null, { toggleMenu })(Navicon);
