import React, { PureComponent } from 'react';
import { FastField } from 'formik';
import classnames from 'classnames';
import capitalize from 'lodash/capitalize';
import Error from 'components/ui/formik/Error';

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
    const { name, label = capitalize(name), noError, ...props } = this.props;

    return (
      <FastField
        name={name}
        render={({ field, form }) => {
          return (
            <>
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
                  {label}
                </span>
                <span
                  className={classnames('textinput-2', {
                    value: field.value,
                    error: form.touched[name] && form.errors[name]
                  })}
                >
                  {label}
                </span>
              </div>
              {!noError && <Error name={name} />}
            </>
          );
        }}
      />
    );
  }
}
