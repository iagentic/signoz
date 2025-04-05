package telemetrylogs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	schema "github.com/SigNoz/signoz-otel-collector/cmd/signozschemamigrator/schema_migrator"
	"github.com/SigNoz/signoz/pkg/types"
	"github.com/huandu/go-sqlbuilder"
)

var (
	mainColumns = map[string]*schema.Column{
		"ts_bucket_start":      {Name: "ts_bucket_start", Type: schema.ColumnTypeUInt64},
		"resource_fingerprint": {Name: "resource_fingerprint", Type: schema.ColumnTypeString},

		"timestamp":          {Name: "timestamp", Type: schema.ColumnTypeUInt64},
		"observed_timestamp": {Name: "observed_timestamp", Type: schema.ColumnTypeUInt64},
		"id":                 {Name: "id", Type: schema.ColumnTypeString},
		"trace_id":           {Name: "trace_id", Type: schema.ColumnTypeString},
		"span_id":            {Name: "span_id", Type: schema.ColumnTypeString},
		"trace_flags":        {Name: "trace_flags", Type: schema.ColumnTypeUInt32},
		"severity_text":      {Name: "severity_text", Type: schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString}},
		"severity_number":    {Name: "severity_number", Type: schema.ColumnTypeUInt8},
		"body":               {Name: "body", Type: schema.ColumnTypeString},
		"attributes_string": {Name: "attributes_string", Type: schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeString,
		}},
		"attributes_number": {Name: "attributes_int", Type: schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeInt64,
		}},
		"attributes_bool": {Name: "attributes_bool", Type: schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeUInt8,
		}},
		"resources_string": {Name: "resources_string", Type: schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeString,
		}},
		"scope_name":    {Name: "scope_name", Type: schema.ColumnTypeString},
		"scope_version": {Name: "scope_version", Type: schema.ColumnTypeString},
		"scope_string": {Name: "scope_string", Type: schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeString,
		}},
	}

	ErrColumnNotFound = errors.New("column not found")
	ErrBetweenValues  = errors.New("(not) between operator requires two values")
	ErrInValues       = errors.New("(not) in operator requires a list of values")
)

var _ types.ConditionBuilder = &conditionBuilder{}

type conditionBuilder struct {
}

func NewConditionBuilder() types.ConditionBuilder {
	return &conditionBuilder{}
}

func (c *conditionBuilder) GetColumn(ctx context.Context, key types.TelemetryFieldKey) (*schema.Column, error) {

	switch key.FieldContext {
	case types.FieldContextResource:
		return mainColumns["resources_string"], nil
	case types.FieldContextScope:
		switch key.Name {
		case "name", "scope.name", "scope_name":
			return mainColumns["scope_name"], nil
		case "version", "scope.version", "scope_version":
			return mainColumns["scope_version"], nil
		}
		return mainColumns["scope_string"], nil
	case types.FieldContextAttribute:
		switch key.FieldDataType {
		case types.FieldDataTypeString:
			return mainColumns["attributes_string"], nil
		case types.FieldDataTypeInt64, types.FieldDataTypeFloat64, types.FieldDataTypeNumber:
			return mainColumns["attributes_number"], nil
		case types.FieldDataTypeBool:
			return mainColumns["attributes_bool"], nil
		}
	case types.FieldContextLog:
		col, ok := mainColumns[key.Name]
		if !ok {
			return nil, ErrColumnNotFound
		}
		return col, nil
	}

	return nil, ErrColumnNotFound
}

func keyToMaterializedColumnName(key types.TelemetryFieldKey) string {
	return fmt.Sprintf("%s_%s_%s", key.FieldContext, key.FieldDataType.String(), strings.ReplaceAll(key.Name, ".", "$$"))
}

func (c *conditionBuilder) getFieldKeyName(ctx context.Context, key types.TelemetryFieldKey) (string, error) {
	column, err := c.GetColumn(ctx, key)
	if err != nil {
		return "", err
	}

	switch column.Type {
	case schema.ColumnTypeString,
		schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
		schema.ColumnTypeUInt64,
		schema.ColumnTypeUInt32,
		schema.ColumnTypeUInt8:
		return column.Name, nil
	case schema.MapColumnType{
		KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
		ValueType: schema.ColumnTypeString,
	}:
		// a key could have been materialized, if so return the materialized column name
		if key.Materialized {
			return keyToMaterializedColumnName(key), nil
		}
		return fmt.Sprintf("%s['%s']", column.Name, key.Name), nil
	case schema.MapColumnType{
		KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
		ValueType: schema.ColumnTypeInt64,
	}:
		// a key could have been materialized, if so return the materialized column name
		if key.Materialized {
			return keyToMaterializedColumnName(key), nil
		}
		return fmt.Sprintf("%s['%s']", column.Name, key.Name), nil
	case schema.MapColumnType{
		KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
		ValueType: schema.ColumnTypeUInt8,
	}:
		// a key could have been materialized, if so return the materialized column name
		if key.Materialized {
			return keyToMaterializedColumnName(key), nil
		}
		return fmt.Sprintf("%s['%s']", column.Name, key.Name), nil
	}
	// should not reach here
	return column.Name, nil
}

func (c *conditionBuilder) GetCondition(
	ctx context.Context,
	key types.TelemetryFieldKey,
	operator types.FilterOperator,
	value any,
	sb *sqlbuilder.SelectBuilder,
) (*sqlbuilder.SelectBuilder, error) {
	column, err := c.GetColumn(ctx, key)
	if err != nil {
		return nil, err
	}

	fieldKeyName, err := c.getFieldKeyName(ctx, key)
	if err != nil {
		return nil, err
	}

	// regular operators
	switch operator {
	// regular operators
	case types.FilterOperatorEqual:
		sb.Where(sb.E(fieldKeyName, value))
	case types.FilterOperatorNotEqual:
		sb.Where(sb.NE(fieldKeyName, value))
	case types.FilterOperatorGreaterThan:
		sb.Where(sb.G(fieldKeyName, value))
	case types.FilterOperatorGreaterThanOrEq:
		sb.Where(sb.GE(fieldKeyName, value))
	case types.FilterOperatorLessThan:
		sb.Where(sb.LT(fieldKeyName, value))
	case types.FilterOperatorLessThanOrEq:
		sb.Where(sb.LE(fieldKeyName, value))

	// like and not like
	case types.FilterOperatorLike:
		sb.Where(sb.Like(fieldKeyName, value))
	case types.FilterOperatorNotLike:
		sb.Where(sb.NotLike(fieldKeyName, value))
	case types.FilterOperatorILike:
		sb.Where(sb.ILike(fieldKeyName, value))
	case types.FilterOperatorNotILike:
		sb.Where(sb.NotILike(fieldKeyName, value))

	// between and not between
	case types.FilterOperatorBetween:
		values, ok := value.([]any)
		if !ok {
			return nil, ErrBetweenValues
		}
		if len(values) != 2 {
			return nil, ErrBetweenValues
		}
		sb.Where(sb.Between(fieldKeyName, values[0], values[1]))
	case types.FilterOperatorNotBetween:
		values, ok := value.([]any)
		if !ok {
			return nil, ErrBetweenValues
		}
		if len(values) != 2 {
			return nil, ErrBetweenValues
		}
		sb.Where(sb.NotBetween(fieldKeyName, values[0], values[1]))

	// in and not in
	case types.FilterOperatorIn:
		values, ok := value.([]any)
		if !ok {
			return nil, ErrInValues
		}
		sb.Where(sb.In(fieldKeyName, values...))
	case types.FilterOperatorNotIn:
		values, ok := value.([]any)
		if !ok {
			return nil, ErrInValues
		}
		sb.Where(sb.NotIn(fieldKeyName, values...))

	// exists and not exists
	// but how could you live and have no story to tell
	// in the UI based query builder, `exists` and `not exists` are used for
	// key membership checks, so depending on the column type, the condition changes
	case types.FilterOperatorExists, types.FilterOperatorNotExists:
		var value any
		switch column.Type {
		case schema.ColumnTypeString, schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString}:
			value = ""
			if operator == types.FilterOperatorExists {
				sb.Where(sb.NE(fieldKeyName, value))
			} else {
				sb.Where(sb.E(fieldKeyName, value))
			}
		case schema.ColumnTypeUInt64, schema.ColumnTypeUInt32, schema.ColumnTypeUInt8:
			value = 0
			if operator == types.FilterOperatorExists {
				sb.Where(sb.NE(fieldKeyName, value))
			} else {
				sb.Where(sb.E(fieldKeyName, value))
			}
		case schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeString,
		}, schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeUInt8,
		}, schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeInt64,
		}:
			leftOperand := fmt.Sprintf("mapContains(%s, '%s')", column.Name, key.Name)
			if operator == types.FilterOperatorExists {
				sb.Where(sb.E(leftOperand, true))
			} else {
				sb.Where(sb.NE(leftOperand, true))
			}
		default:
			return nil, fmt.Errorf("exists operator is not supported for column type %s", column.Type)
		}
	}
	return sb, nil
}
