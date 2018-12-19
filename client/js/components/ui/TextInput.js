import React, { PureComponent } from 'react';
import { FastField } from 'formik';
import classnames from 'classnames';
import capitalize from 'lodash/capitalize';
import Error from 'components/ui/formik/Error';

const getValue = (e, trim) => {
  let v = e.target.value;

  if (trim) {
    v = v.trim();
  }

  if (e.target.type === 'number') {
    v = parseFloat(v);
    /* eslint-disable-next-line no-self-compare */
    if (v !== v) {
      v = '';
    }
  }

  return v;
};

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
                autoCapitalize="off"
                autoCorrect="off"
                autoComplete="off"
                spellCheck="false"
                ref={this.input}
                onFocus={this.handleFocus}
                {...field}
                {...props}
                onChange={e => {
                  let v = getValue(e, !noTrim);
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
                  if (blurTransform) {
                    const v = blurTransform(getValue(e));

                    if (v && v !== field.value) {
                      form.setFieldValue(name, v, false);
                    }
                  }

                  field.onBlur(e);
                  if (props.onBlur) {
                    props.onBlur(e);
                  }
                }}
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
        )}
      />
    );
  }
}
