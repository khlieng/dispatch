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
    const {
      name,
      label = capitalize(name),
      noError,
      noTrim,
      transform,
      blurTransform,
      ...props
    } = this.props;

    return (
      <FastField
        name={name}
        render={({ field, form }) => (
          <>
            <div className="textinput">
              <input
                className={field.value && 'value'}
                type="text"
                name={name}
                id={name}
                autoCapitalize="off"
                autoCorrect="off"
                autoComplete="off"
                spellCheck="false"
                ref={this.input}
                onFocus={this.handleFocus}
                {...field}
                {...props}
                onChange={e => {
                  let v = e.target.value;

                  if (!noTrim) {
                    v = v.trim();
                  }

                  if (transform) {
                    v = transform(v);
                  }

                  if (v !== field.value) {
                    form.setFieldValue(name, v);

                    if (props.onChange) {
                      props.onChange(e);
                    }
                  }
                }}
                onBlur={e => {
                  field.onBlur(e);
                  if (props.onBlur) {
                    props.onBlur(e);
                  }

                  if (blurTransform) {
                    const v = blurTransform(e.target.value);

                    if (v && v !== field.value) {
                      form.setFieldValue(name, v);
                    }
                  }
                }}
              />
              <label
                htmlFor={name}
                className={classnames('textinput-label', 'textinput-1', {
                  value: field.value,
                  error: form.touched[name] && form.errors[name]
                })}
              >
                {label}
              </label>
              <span
                className={classnames('textinput-label', 'textinput-2', {
                  value: field.value,
                  error: form.touched[name] && form.errors[name]
                })}
              >
                {label}
              </span>
            </div>
            {!noError && <Error name={name} />}
          </>
        )}
      />
    );
  }
}
