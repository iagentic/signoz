package telemetryspans

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
		"trace_id":           {Name: "trace_id", Type: schema.ColumnTypeString},
		"span_id":            {Name: "span_id", Type: schema.ColumnTypeString},
		"trace_state":        {Name: "trace_state", Type: schema.ColumnTypeString},
		"parent_span_id":     {Name: "parent_span_id", Type: schema.ColumnTypeString},
		"flags":              {Name: "flags", Type: schema.ColumnTypeUInt32},
		"name":               {Name: "name", Type: schema.ColumnTypeString},
		"kind":               {Name: "kind", Type: schema.ColumnTypeInt8},
		"kind_string":        {Name: "kind_string", Type: schema.ColumnTypeString},
		"duration_nano":      {Name: "duration_nano", Type: schema.ColumnTypeUInt64},
		"status_code":        {Name: "status_code", Type: schema.ColumnTypeInt16},
		"status_message":     {Name: "status_message", Type: schema.ColumnTypeString},
		"status_code_string": {Name: "status_code_string", Type: schema.ColumnTypeString},

		"attributes_string": {Name: "attributes_string", Type: schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeString,
		}},
		"attributes_number": {Name: "attributes_number", Type: schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeFloat64,
		}},
		"attributes_bool": {Name: "attributes_bool", Type: schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeUInt8,
		}},
		"resources_string": {Name: "resources_string", Type: schema.MapColumnType{
			KeyType:   schema.LowCardinalityColumnType{ElementType: schema.ColumnTypeString},
			ValueType: schema.ColumnTypeString,
		}},
		"events": {Name: "events", Type: schema.ArrayColumnType{
			ElementType: schema.ColumnTypeString,
		}},
		"links":                                 {Name: "links", Type: schema.ColumnTypeString},
		"response_status_code":                  {Name: "response_status_code", Type: schema.ColumnTypeString},
		"external_http_url":                     {Name: "external_http_url", Type: schema.ColumnTypeString},
		"http_url":                              {Name: "http_url", Type: schema.ColumnTypeString},
		"external_http_method":                  {Name: "external_http_method", Type: schema.ColumnTypeString},
		"http_method":                           {Name: "http_method", Type: schema.ColumnTypeString},
		"http_host":                             {Name: "http_host", Type: schema.ColumnTypeString},
		"db_name":                               {Name: "db_name", Type: schema.ColumnTypeString},
		"db_operation":                          {Name: "db_operation", Type: schema.ColumnTypeString},
		"has_error":                             {Name: "has_error", Type: schema.ColumnTypeUInt8},
		"is_remote":                             {Name: "is_remote", Type: schema.ColumnTypeString},
		"resource_string_service$$name":         {Name: "resource_string_service$$name", Type: schema.ColumnTypeString},
		"attribute_string_http$$route":          {Name: "attribute_string_http$$route", Type: schema.ColumnTypeString},
		"attribute_string_messaging$$system":    {Name: "attribute_string_messaging$$system", Type: schema.ColumnTypeString},
		"attribute_string_messaging$$operation": {Name: "attribute_string_messaging$$operation", Type: schema.ColumnTypeString},
		"attribute_string_db$$system":           {Name: "attribute_string_db$$system", Type: schema.ColumnTypeString},
		"attribute_string_rpc$$system":          {Name: "attribute_string_rpc$$system", Type: schema.ColumnTypeString},
		"attribute_string_rpc$$service":         {Name: "attribute_string_rpc$$service", Type: schema.ColumnTypeString},
		"attribute_string_rpc$$method":          {Name: "attribute_string_rpc$$method", Type: schema.ColumnTypeString},
		"attribute_string_peer$$service":        {Name: "attribute_string_peer$$service", Type: schema.ColumnTypeString},

		"traceID":          {Name: "traceID", Type: schema.ColumnTypeString},
		"spanID":           {Name: "spanID", Type: schema.ColumnTypeString},
		"parentSpanID":     {Name: "parentSpanID", Type: schema.ColumnTypeString},
		"spanKind":         {Name: "spanKind", Type: schema.ColumnTypeString},
		"durationNano":     {Name: "durationNano", Type: schema.ColumnTypeUInt64},
		"statusCode":       {Name: "statusCode", Type: schema.ColumnTypeInt16},
		"statusMessage":    {Name: "statusMessage", Type: schema.ColumnTypeString},
		"statusCodeString": {Name: "statusCodeString", Type: schema.ColumnTypeString},

		"references":         {Name: "references", Type: schema.ColumnTypeString},
		"responseStatusCode": {Name: "responseStatusCode", Type: schema.ColumnTypeString},
		"externalHttpUrl":    {Name: "externalHttpUrl", Type: schema.ColumnTypeString},
		"httpUrl":            {Name: "httpUrl", Type: schema.ColumnTypeString},
		"externalHttpMethod": {Name: "externalHttpMethod", Type: schema.ColumnTypeString},

		"resource_string_service$$name_exists":         {Name: "resource_string_service$$name_exists", Type: schema.ColumnTypeUInt8},
		"attribute_string_http$$route_exists":          {Name: "attribute_string_http$$route_exists", Type: schema.ColumnTypeUInt8},
		"attribute_string_messaging$$system_exists":    {Name: "attribute_string_messaging$$system_exists", Type: schema.ColumnTypeUInt8},
		"attribute_string_messaging$$operation_exists": {Name: "attribute_string_messaging$$operation_exists", Type: schema.ColumnTypeUInt8},
		"attribute_string_db$$system_exists":           {Name: "attribute_string_db$$system_exists", Type: schema.ColumnTypeUInt8},
		"attribute_string_rpc$$system_exists":          {Name: "attribute_string_rpc$$system_exists", Type: schema.ColumnTypeUInt8},
		"attribute_string_rpc$$service_exists":         {Name: "attribute_string_rpc$$service_exists", Type: schema.ColumnTypeUInt8},
		"attribute_string_rpc$$method_exists":          {Name: "attribute_string_rpc$$method_exists", Type: schema.ColumnTypeUInt8},
		"attribute_string_peer$$service_exists":        {Name: "attribute_string_peer$$service_exists", Type: schema.ColumnTypeUInt8},
	}

	ErrColumnNotFound = errors.New("column not found")
	ErrBetweenValues  = errors.New("(not) between operator requires two values")
	ErrInValues       = errors.New("(not) in operator requires a list of values")
)

// interface check
var _ types.ConditionBuilder = &conditionBuilder{}

type conditionBuilder struct {
}

func NewConditionBuilder() types.ConditionBuilder {
	return &conditionBuilder{}
}

func keyToMaterializedColumnName(key types.TelemetryFieldKey) string {
	return fmt.Sprintf("%s_%s_%s", key.FieldContext, key.FieldDataType.String(), strings.ReplaceAll(key.Name, ".", "$$"))
}

func (c *conditionBuilder) GetColumn(ctx context.Context, key types.TelemetryFieldKey) (*schema.Column, error) {

	switch key.FieldContext {
	case types.FieldContextResource:
		return mainColumns["resources_string"], nil
	case types.FieldContextScope:
		return nil, ErrColumnNotFound
	case types.FieldContextAttribute:
		switch key.FieldDataType {
		case types.FieldDataTypeString:
			return mainColumns["attributes_string"], nil
		case types.FieldDataTypeInt64, types.FieldDataTypeFloat64, types.FieldDataTypeNumber:
			return mainColumns["attributes_number"], nil
		case types.FieldDataTypeBool:
			return mainColumns["attributes_bool"], nil
		}
	case types.FieldContextSpan:
		col, ok := mainColumns[key.Name]
		if !ok {
			return nil, ErrColumnNotFound
		}
		return col, nil
	}

	return nil, ErrColumnNotFound
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
