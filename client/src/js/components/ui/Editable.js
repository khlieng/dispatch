import React, { PureComponent } from 'react';
import { stringWidth } from 'utils';

export default class Editable extends PureComponent {
  static defaultProps = {
    editable: true
  };

  state = {
    editing: false
  };

  componentWillReceiveProps(nextProps) {
    if (this.state.editing && nextProps.value !== this.props.value) {
      this.setState({
        width: this.getInputWidth(nextProps.value)
      });
    }
  }

  componentDidUpdate(prevProps, prevState) {
    if (!prevState.editing && this.state.editing) {
      // eslint-disable-next-line react/no-did-update-set-state
      this.setState({
        width: this.getInputWidth(this.props.value)
      });
    }
  }

  getInputWidth(value) {
    if (this.input) {
      const style = window.getComputedStyle(this.input);
      const padding = parseInt(style.paddingLeft, 10) + parseInt(style.paddingRight, 10);
      // Make sure the width is atleast 1px so the caret always shows
      const width = stringWidth(value, style.font) || 1;
      return padding + width;
    }
  }

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

  handleBlur = e => {
    const { onBlur } = this.props;
    this.stopEditing();
    if (onBlur) {
      onBlur(e.target.value);
    }
  };

  handleChange = e => this.props.onChange(e.target.value);

  handleKey = e => {
    if (e.key === 'Enter') {
      this.handleBlur(e);
    }
  };

  inputRef = el => { this.input = el; }

  render() {
    const { children, className, value } = this.props;

    const style = {
      width: this.state.width
    };

    return (
      this.state.editing ?
        <input
          autoFocus
          ref={this.inputRef}
          className={className}
          type="text"
          value={value}
          onBlur={this.handleBlur}
          onChange={this.handleChange}
          onKeyDown={this.handleKey}
          style={style}
          spellCheck={false}
        /> :
        <div onClick={this.startEditing}>{children}</div>
    );
  }
}
