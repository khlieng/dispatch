import React, { PureComponent } from 'react';

const style = {
  background: 'none',
  font: 'inherit'
};

export default class Editable extends PureComponent {
  state = { editing: false };

  startEditing = () => {
    if (this.props.editable) {
      this.initialValue = this.props.value;
      this.setState({ editing: true });
    }
  };

  stopEditing = () => {
    const { validate, value, onChange } = this.props;
    if (validate && !validate(value)) {
      onChange(this.initialValue);
    }
    this.setState({ editing: false });
  };

  handleKey = e => {
    if (e.key === 'Enter') {
      this.stopEditing();
    }
  };

  handleChange = e => this.props.onChange(e.target.value);

  render() {
    const { children, className, value } = this.props;
    return (
      <div onClick={this.startEditing}>
        {this.state.editing ?
          <input
            autoFocus
            className={className}
            style={style}
            type="text"
            value={value}
            onBlur={this.stopEditing}
            onChange={this.handleChange}
            onKeyDown={this.handleKey}
          /> :
          children
        }
      </div>
    );
  }
}
