import React from 'react';
import { Field } from 'formik';
import Checkbox from 'components/ui/Checkbox';

const FormikCheckbox = ({ name, onChange, ...props }) => (
  <Field
    name={name}
    render={({ field, form }) => (
      <Checkbox
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
    )}
  />
);

export default FormikCheckbox;
