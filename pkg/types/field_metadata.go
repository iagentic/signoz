package types

// FieldMetadataService defines the interface for retrieving field metadata
type FieldMetadataService interface {
	// GetFieldMetadata retrieves metadata for a specific field
	GetFieldMetadata(key string, signalType SignalType) (*Field, error)

	// GetFieldsMetadata retrieves metadata for multiple fields
	GetFieldsMetadata(keys []string, signalType SignalType) ([]*Field, error)

	// GetAllFieldsMetadata retrieves metadata for all available fields for a signal type
	GetAllFieldsMetadata(signalType SignalType) ([]*Field, error)

	// GetSampleValues retrieves sample values for a specific field
	GetSampleValues(fieldKey *FieldKey, signalType SignalType, limit int) ([]any, error)

	// GetRelatedValues retrieves values related to a specific field
	GetRelatedValues(fieldKey *FieldKey, value any, signalType SignalType, limit int) ([]any, error)
}
