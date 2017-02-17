import React, { PureComponent } from 'react';

export default class FileInput extends PureComponent {
  componentWillMount() {
    this.input = window.document.createElement('input');
    this.input.setAttribute('type', 'file');

    this.input.addEventListener('change', e => {
      const file = e.target.files[0];
      const reader = new FileReader();

      reader.onload = () => {
        console.log(reader.result.byteLength);
        this.props.onChange(file.name, reader.result);
      };

      reader.readAsArrayBuffer(file);
    });
  }

  handleClick = () => this.input.click();

  render() {
    return (
      <button className="input-file" onClick={this.handleClick}>{this.props.name}</button>
    );
  }
}
