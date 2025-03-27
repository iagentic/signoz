package types

import "context"

type FieldContext string

const (
	FieldContextMetric    FieldContext = "metric"
	FieldContextLog       FieldContext = "log"
	FieldContextSpan      FieldContext = "span"
	FieldContextTrace     FieldContext = "trace"
	FieldContextResource  FieldContext = "resource"
	FieldContextScope     FieldContext = "scope"
	FieldContextAttribute FieldContext = "attribute"
	FieldContextEvent     FieldContext = "event"
	FieldContextAll       FieldContext = "all"
)

func (f FieldContext) String() string {
	return string(f)
}

func FieldContextFromString(s string) FieldContext {
	switch s {
	case "resource":
		return FieldContextResource
	case "scope":
		return FieldContextScope
	case "tag", "attribute":
		return FieldContextAttribute
	case "event":
		return FieldContextEvent
	case "spanfield", "span":
		return FieldContextSpan
	case "logfield", "log":
		return FieldContextLog
	case "metric":
		return FieldContextMetric
	default:
		return FieldContextAll
	}
}

type FieldDataType string

const (
	FieldDataTypeString FieldDataType = "string"
	FieldDataTypeBool   FieldDataType = "bool"
	FieldDataTypeInt    FieldDataType = "int"
	FieldDataTypeFloat  FieldDataType = "float"
	FieldDataTypeNumber FieldDataType = "number"
	FieldDataTypeAll    FieldDataType = "all"
)

func (f FieldDataType) String() string {
	return string(f)
}

func FieldDataTypeFromString(s string) FieldDataType {
	switch s {
	case "string":
		return FieldDataTypeString
	case "bool":
		return FieldDataTypeBool
	case "int":
		return FieldDataTypeInt
	case "float":
		return FieldDataTypeFloat
	case "number":
		return FieldDataTypeNumber
	default:
		return FieldDataTypeAll
	}
}

type FieldKeySelectorType string

const (
	FieldKeySelectorTypeExact FieldKeySelectorType = "exact"
	FieldKeySelectorTypeFuzzy FieldKeySelectorType = "fuzzy"
)

type TelemetryFieldKey struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Unit          string        `json:"unit"`
	FieldContext  FieldContext  `json:"fieldContext"`
	FieldDataType FieldDataType `json:"fieldDataType"`
	Materialized  bool          `json:"-"`
}

type ExistingFieldSelection struct {
	Key   TelemetryFieldKey `json:"key"`
	Value any               `json:"value"`
}

type FieldKeySelector struct {
	FieldContext  FieldContext         `json:"fieldContext"`
	FieldDataType FieldDataType        `json:"fieldDataType"`
	Name          string               `json:"name"`
	SelectorType  FieldKeySelectorType `json:"selectorType"`
	Limit         int                  `json:"limit"`
}

type Metadata interface {
	// GetKeys returns a map of field keys by name, there can be multiple keys with the same name
	// if they have different types or data types.
	GetKeys(ctx context.Context, fieldKeySelector FieldKeySelector) (map[string][]TelemetryFieldKey, error)

	// GetKey returns a list of keys with the given name.
	GetKey(ctx context.Context, fieldKeySelector FieldKeySelector) ([]TelemetryFieldKey, error)

	// GetRelatedValues returns a list of related values for the given key name
	// and the existing selection of keys.
	GetRelatedValues(ctx context.Context, fieldKeySelector FieldKeySelector, existingSelections []ExistingFieldSelection) (any, error)

	// GetAllValues returns a list of all values.
	GetAllValues(ctx context.Context, fieldKeySelector FieldKeySelector) (any, error)
}
