import React, { PureComponent } from 'react';
import Button from 'components/ui/Button';

export default class FileInput extends PureComponent {
  static defaultProps = {
    type: 'text'
  };

  constructor(props) {
    super(props);

    this.input = window.document.createElement('input');
    this.input.setAttribute('type', 'file');

    this.input.addEventListener('change', e => {
      const file = e.target.files[0];
      const reader = new FileReader();
      const { onChange, type } = this.props;

      reader.onload = () => {
        onChange(file.name, reader.result);
      };

      switch (type) {
        case 'binary':
          reader.readAsArrayBuffer(file);
          break;

        case 'text':
          reader.readAsText(file);
          break;

        default:
          reader.readAsText(file);
      }
    });
  }

  handleClick = () => this.input.click();

  render() {
    return (
      <Button className="input-file" onClick={this.handleClick}>
        {this.props.name}
      </Button>
    );
  }
}
