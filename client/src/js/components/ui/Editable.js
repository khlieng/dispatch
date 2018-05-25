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
      this.updateInputWidth(nextProps.value);
    }
  }

  componentDidUpdate(prevProps, prevState) {
    if (!prevState.editing && this.state.editing) {
      // eslint-disable-next-line react/no-did-update-set-state
      this.updateInputWidth(this.props.value);
    }
  }

  updateInputWidth = value => {
    if (this.input) {
      const style = window.getComputedStyle(this.input);
      const padding = parseInt(style.paddingRight, 10);
      // Make sure the width is atleast 1px so the caret always shows
      const width =
        stringWidth(value, `${style.fontSize} ${style.fontFamily}`) || 1;

      this.setState({
        width: width + padding * 2,
        indent: padding
      });
    }
  };

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

  handleFocus = e => {
    const val = e.target.value;
    e.target.value = '';
    e.target.value = val;
  };

  inputRef = el => {
    this.input = el;
  };

  render() {
    const { children, className, value } = this.props;

    const style = {
      width: this.state.width,
      textIndent: this.state.indent,
      paddingLeft: 0
    };

    return this.state.editing ? (
      <input
        autoFocus
        ref={this.inputRef}
        className={className}
        type="text"
        value={value}
        onBlur={this.handleBlur}
        onChange={this.handleChange}
        onKeyDown={this.handleKey}
        onFocus={this.handleFocus}
        style={style}
        spellCheck={false}
      />
    ) : (
      <div onClick={this.startEditing}>{children}</div>
    );
  }
}
