import React, { PureComponent, createRef } from 'react';
import cn from 'classnames';
import { stringWidth } from 'utils';

export default class Editable extends PureComponent {
  static defaultProps = {
    editable: true
  };

  inputEl = createRef();

  state = {
    editing: false
  };

  componentDidUpdate(prevProps, prevState) {
    if (!prevState.editing && this.state.editing) {
      // eslint-disable-next-line react/no-did-update-set-state
      this.updateInputWidth(this.props.value);
      this.inputEl.current.focus();
    } else if (this.state.editing && prevProps.value !== this.props.value) {
      this.updateInputWidth(this.props.value);
    }
  }

  updateInputWidth = value => {
    if (this.inputEl.current) {
      const style = window.getComputedStyle(this.inputEl.current);
      const padding = parseInt(style.paddingRight, 10);
      // Make sure the width is at least 1px so the caret always shows
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

  render() {
    const { children, className, editable, value } = this.props;

    const style = {
      width: this.state.width,
      textIndent: this.state.indent,
      paddingLeft: 0
    };

    return this.state.editing ? (
      <input
        ref={this.inputEl}
        className={`editable-wrap ${className}`}
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
      <div
        className={cn('editable-wrap', {
          'editable-wrap-editable': editable
        })}
        onClick={this.startEditing}
      >
        {children}
      </div>
    );
  }
}
