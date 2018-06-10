import React from 'react';
import { Field } from 'formik';

const Checkbox = ({ name, label, onChange, ...props }) => (
  <Field
    name={name}
    render={({ field, form }) => (
      <label htmlFor={name}>
        {label && <div>{label}</div>}
        <input
          type="checkbox"
          id={name}
          name={name}
          checked={field.value}
          onChange={e => {
            form.setFieldTouched(name, true);
            field.onChange(e);
            if (onChange) {
              onChange(e);
            }
          }}
          {...props}
        />
        <span />
      </label>
    )}
  />
);

export default Checkbox;
