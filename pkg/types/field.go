package types

import "fmt"

// FieldContext is the context of a field
type FieldContext string

const (
	// FieldContextUnknown is used when the context is unknown
	FieldContextUnknown FieldContext = "unknown"
	// FieldContextAttribute is used for attributes of {span, metric, log}
	FieldContextAttribute FieldContext = "attribute"
	// FieldContextResource is used for resource attributes of {span, metric, log}
	FieldContextResource FieldContext = "resource"
	// FieldContextPrimary is used for primary fields of {span, metric, log}
	// Example: name, kind, severity_text, span_id, trace_id, etc.
	FieldContextPrimary FieldContext = "primary"
)

func (q FieldContext) Validate() error {
	switch q {
	case FieldContextAttribute, FieldContextResource, FieldContextPrimary:
		return nil
	default:
		return fmt.Errorf("invalid field context: %s, valid contexts are %v", q, []FieldContext{FieldContextAttribute, FieldContextResource, FieldContextPrimary})
	}
}

// FieldDataType is the data type of a field
type FieldDataType string

const (
	// FieldDataTypeUnknown is used when the data type is unknown
	FieldDataTypeUnknown FieldDataType = "unknown"
	// FieldDataTypeString is used for string fields
	FieldDataTypeString FieldDataType = "string"
	// FieldDataTypeNumber is used for number fields
	FieldDataTypeNumber FieldDataType = "number"
	// FieldDataTypeBoolean is used for boolean fields
	FieldDataTypeBoolean FieldDataType = "boolean"
	// FieldDataTypeStringArray is used for array of strings fields
	FieldDataTypeStringArray FieldDataType = "[]string"
	// FieldDataTypeNumberArray is used for array of numbers fields
	FieldDataTypeNumberArray FieldDataType = "[]number"
	// FieldDataTypeBooleanArray is used for array of booleans fields
	FieldDataTypeBooleanArray FieldDataType = "[]boolean"
	// FieldDataTypeObjectArray is used for array of any fields
	FieldDataTypeObjectArray FieldDataType = "[]any"
)

func (q FieldDataType) Validate() error {
	switch q {
	case FieldDataTypeString, FieldDataTypeNumber, FieldDataTypeBoolean, FieldDataTypeStringArray, FieldDataTypeNumberArray, FieldDataTypeBooleanArray, FieldDataTypeObjectArray:
		return nil
	default:
		return fmt.Errorf("invalid field data type: %s, valid data types are %v", q, []FieldDataType{FieldDataTypeString, FieldDataTypeNumber, FieldDataTypeBoolean, FieldDataTypeStringArray, FieldDataTypeNumberArray, FieldDataTypeBooleanArray, FieldDataTypeObjectArray})
	}
}

// FieldKey represents a field in a signal
type FieldKey struct {
	// Context is the context of the field
	Context FieldContext `json:"context"`
	// DataType is the data type of the field
	DataType FieldDataType `json:"dataType"`
	// SignalType is the signal type of the field
	SignalType SignalType `json:"signalType"`
	// Key is the key of the field
	Key string `json:"key"`
}

func (q *FieldKey) Validate() error {
	if err := q.Context.Validate(); err != nil {
		return err
	}
	if err := q.DataType.Validate(); err != nil {
		return err
	}
	if err := q.SignalType.Validate(); err != nil {
		return err
	}
	return nil
}

type Field struct {
	FieldKey
	// SampleValues contains example values for this field
	SampleValues []interface{} `json:"sampleValues,omitempty"`
	// RelatedValues contains values related to this field that might be useful
	RelatedValues []interface{} `json:"relatedValues,omitempty"`
	// Cardinality represents the approximate number of unique values for this field
	Cardinality int64 `json:"cardinality,omitempty"`
}

func (q *Field) Validate() error {
	if err := q.FieldKey.Validate(); err != nil {
		return err
	}
	return nil
}
