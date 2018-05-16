import React from 'react';
import { Field } from 'formik';

const TextInput = ({ name, placeholder, ...props }) => (
  <Field
    name={name}
    render={({ field }) => (
      <div className="textinput">
        <input
          className={field.value ? 'value' : null}
          type="text"
          name={name}
          {...field}
          {...props}
        />
        <span className={field.value ? 'textinput-1 value' : 'textinput-1'}>
          {placeholder}
        </span>
        <span className={field.value ? 'textinput-2 value' : 'textinput-2'}>
          {placeholder}
        </span>
      </div>
    )}
  />
);

export default TextInput;
