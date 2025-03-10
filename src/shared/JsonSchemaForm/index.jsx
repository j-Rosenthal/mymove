import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { isEmpty, uniq } from 'lodash';
import { Field, reduxForm } from 'redux-form';

import SchemaField, { ALWAYS_REQUIRED_KEY } from './JsonSchemaField';

import 'shared/JsonSchemaForm/index.css';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';

const renderGroupOrField = (fieldName, fields, uiSchema, nameSpace) => {
  /* TODO:
   dates look wonky in chrome
   styling in accordance with USWDS
   validate group names don't colide with field names
  */
  const group = uiSchema.groups && uiSchema.groups[fieldName];
  const isRef = fields[fieldName] && fields[fieldName].$$ref && fields[fieldName].properties;
  const isCustom = uiSchema.custom_components && uiSchema.custom_components[fieldName];
  if (group) {
    const keys = group.fields;
    return (
      <fieldset className="usa-fieldset" key={fieldName}>
        <p htmlFor={fieldName}>{group.title}</p>
        {keys.map((f) => renderGroupOrField(f, fields, uiSchema, nameSpace))}
      </fieldset>
    );
  }
  if (isCustom) {
    return (
      <Fragment key={fieldName}>
        <p>{fields[fieldName].title}</p>
        <Field name={fieldName} component={uiSchema.custom_components[fieldName]} />
      </Fragment>
    );
  }
  if (isRef) {
    const refName = fields[fieldName].$$ref.split('/').pop();
    const refSchema = uiSchema.definitions[refName];
    return renderSchema(fields[fieldName], refSchema, fieldName);
  }
  return renderField(fieldName, fields, nameSpace);
};

export const renderField = (fieldName, fields, nameSpace) => {
  const field = fields[fieldName];
  if (!field) {
    return undefined;
  }
  return SchemaField.createSchemaField(fieldName, field, nameSpace);
};

// Because we have nested objects it's possible to have
// An object that is not-required that itself has required properties. This makes sense, in that
// If the entire object is omitted (say, an address) then the form is valid, but if a
// single property of the object is included, then all its required properties must be
// as well.
// Therefore, the rules for wether or not a field is required are:
// 1. If it is listed in the top level definition, it's required.
// 2. If it is required and it is an object, its required fields are required
// 3. If it is an object and some value in it has been set, then all it's required fields must be set too
// This is a recusive definition.
export const recursivelyValidateRequiredFields = (values, spec) => {
  const requiredErrors = {};
  // first, check that all required fields are present
  if (spec.required) {
    spec.required.forEach((requiredFieldName) => {
      if (values[requiredFieldName] === undefined || values[requiredFieldName] === '') {
        // check if the required thing is a object, in that case put it on its required fields. Otherwise recurse.
        const schemaForKey = spec.properties[requiredFieldName];
        if (schemaForKey) {
          if (schemaForKey.type === 'object') {
            const subErrors = recursivelyValidateRequiredFields({}, schemaForKey);
            if (!isEmpty(subErrors)) {
              requiredErrors[requiredFieldName] = subErrors;
            }
          } else {
            requiredErrors[requiredFieldName] = 'Required.';
          }
        } else {
          milmoveLog(MILMOVE_LOG_LEVEL.ERROR, 'The schema should have all required fields in it.');
        }
      }
    });
  }

  // now go through every existing value, if its an object, we must recurse to see if its required properties are there.
  Object.keys(values).forEach(function (key) {
    const schemaForKey = spec.properties[key];
    if (schemaForKey) {
      if (schemaForKey.type === 'object') {
        const subErrors = recursivelyValidateRequiredFields(values[key], schemaForKey);
        if (!isEmpty(subErrors)) {
          requiredErrors[key] = subErrors;
        }
      }
    } else {
      milmoveLog(MILMOVE_LOG_LEVEL.ERROR, `The schema should have fields for all present values. Missing ${key}`);
    }
  });

  return requiredErrors;
};

// To validate that fields are required, we look at the list of top level required
// fields and then validate them and their children.
export const validateRequiredFields = (values, form) => {
  const swaggerSpec = form.schema;
  let requiredErrors;
  if (swaggerSpec && !isEmpty(swaggerSpec)) {
    requiredErrors = recursivelyValidateRequiredFields(values, swaggerSpec);
  }
  return requiredErrors;
};

export const validateAdditionalFields = (additionalFields) => {
  return (values, form) => {
    const errors = {};
    additionalFields.forEach((fieldName) => {
      if (values[fieldName] === undefined || values[fieldName] === '' || values[fieldName] === null) {
        errors[fieldName] = 'Required.';
      }
    });

    return errors;
  };
};

// Always Required Fields are fields that are marked as required in swagger, and if they are objects, their sub-required fields.
// Fields like Addresses may not be required, so even though they have required subfields they are not annotated.
export const recursivelyAnnotateRequiredFields = (schema) => {
  if (schema.required) {
    schema.required.forEach((requiredFieldName) => {
      // check if the required thing is a object, in that case put it on its required fields. Otherwise recurse.
      const schemaForKey = schema.properties[requiredFieldName];
      if (schemaForKey) {
        if (schemaForKey.type === 'object') {
          recursivelyAnnotateRequiredFields(schemaForKey);
        } else {
          schemaForKey[ALWAYS_REQUIRED_KEY] = true;
        }
      } else {
        milmoveLog(MILMOVE_LOG_LEVEL.ERROR, 'The schema should have all required fields in it.');
      }
    });
  }
};

export const renderSchema = (schema, uiSchema, nameSpace = '') => {
  if (schema && !isEmpty(schema)) {
    recursivelyAnnotateRequiredFields(schema);

    const fields = schema.properties || {};
    return uiSchema.order.map((i) => renderGroupOrField(i, fields, uiSchema, nameSpace));
  }
  return undefined;
};

export const addUiSchemaRequiredFields = (schema, uiSchema) => {
  if (!uiSchema.requiredFields) return;
  if (!schema.properties) return;
  if (!schema.required) schema.required = uiSchema.requiredFields;
  schema.required = uniq(schema.required.concat(uiSchema.requiredFields));
};

export const JsonSchemaFormBody = (props) => {
  const { schema, uiSchema } = props;

  addUiSchemaRequiredFields(schema, uiSchema);
  const title = uiSchema.title || (schema ? schema.title : '');
  const { description } = uiSchema;
  const { todos } = uiSchema;

  return (
    <>
      <h1>{title}</h1>
      {description && <p>{description}</p>}
      {renderSchema(schema, uiSchema)}
      {todos && (
        <div className="Todo">
          <h3>Todo:</h3>
          {todos}
        </div>
      )}
    </>
  );
};

JsonSchemaFormBody.propTypes = {
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
};

JsonSchemaFormBody.defaultProps = {
  className: 'default',
};

const JsonSchemaForm = (props) => {
  const { className } = props;
  const { handleSubmit, schema, uiSchema } = props;
  return (
    <form className={className} onSubmit={handleSubmit}>
      <JsonSchemaFormBody schema={schema} uiSchema={uiSchema} />
    </form>
  );
};

JsonSchemaForm.propTypes = {
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func.isRequired,
};

JsonSchemaForm.defaultProps = {
  className: 'default',
};

export const reduxifyForm = (name) => reduxForm({ form: name, validate: validateRequiredFields })(JsonSchemaForm);
