import React, { PureComponent } from 'react';
import { Field } from 'formik';
import classnames from 'classnames';

export default class TextInput extends PureComponent {
  constructor(props) {
    super(props);
    this.input = React.createRef();
    window.addEventListener('resize', this.handleResize);
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.handleResize);
  }

  handleResize = () => {
    if (this.scroll) {
      this.scroll = false;
      this.scrollIntoView();
    }
  };

  handleFocus = () => {
    this.scroll = true;
    setTimeout(() => {
      this.scroll = false;
    }, 2000);
  };

  scrollIntoView = () => {
    if (this.input.current.scrollIntoViewIfNeeded) {
      this.input.current.scrollIntoViewIfNeeded();
    } else {
      this.input.current.scrollIntoView();
    }
  };

  render() {
    const { name, placeholder, ...props } = this.props;

    return (
      <Field
        name={name}
        render={({ field, form }) => (
          <div className="textinput">
            <input
              className={field.value && 'value'}
              type="text"
              name={name}
              autoCapitalize="off"
              autoCorrect="off"
              autoComplete="off"
              spellCheck="false"
              ref={this.input}
              onFocus={this.handleFocus}
              {...field}
              {...props}
            />
            <span
              className={classnames('textinput-1', {
                value: field.value,
                error: form.touched[name] && form.errors[name]
              })}
            >
              {placeholder}
            </span>
            <span
              className={classnames('textinput-2', {
                value: field.value,
                error: form.touched[name] && form.errors[name]
              })}
            >
              {placeholder}
            </span>
          </div>
        )}
      />
    );
  }
}
