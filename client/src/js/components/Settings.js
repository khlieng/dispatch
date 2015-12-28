import React, { Component } from 'react';
import pure from 'pure-render-decorator';
import Navicon from './Navicon';

@pure
export default class Settings extends Component {
  render() {
    return (
      <div>
        <Navicon />
      </div>
    );
  }
}
