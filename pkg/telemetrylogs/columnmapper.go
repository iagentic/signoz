package telemetrylogs

import (
	"context"
	"errors"

	schema "github.com/SigNoz/signoz-otel-collector/cmd/signozschemamigrator/schema_migrator"
	"github.com/SigNoz/signoz/pkg/types"
)

var (
	mainColumns = map[string]schema.Column{
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
)

type columnMapper struct {
}

func NewColumnMapper() types.KeyToColumnMapper {
	return &columnMapper{}
}

func (c *columnMapper) GetColumn(ctx context.Context, key types.TelemetryFieldKey) (schema.Column, error) {

	switch key.FieldContext {
	case types.FieldContextResource:
		return mainColumns["resources_string"], nil
	case types.FieldContextScope:
		switch key.Name {
		case "name", "scope.name":
			return mainColumns["scope_name"], nil
		case "version", "scope.version":
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
			return schema.Column{}, ErrColumnNotFound
		}
		return col, nil
	}

	return schema.Column{}, ErrColumnNotFound
}
