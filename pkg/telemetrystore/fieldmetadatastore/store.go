package fieldmetadatastore

import (
	"go.signoz.io/signoz/pkg/telemetrystore"
	"go.signoz.io/signoz/pkg/types"
)

type clickhouseFieldMetadataStore struct {
	store telemetrystore.TelemetryStore
}

// impl check
var _ types.FieldMetadataService = &clickhouseFieldMetadataStore{}

func NewClickhouseFieldMetadataStore(store telemetrystore.TelemetryStore) *clickhouseFieldMetadataStore {
	return &clickhouseFieldMetadataStore{store: store}
}

func (s *clickhouseFieldMetadataStore) GetFieldMetadata(key string, signalType types.SignalType) (*types.Field, error) {
	return nil, nil
}

func (s *clickhouseFieldMetadataStore) GetFieldsMetadata(keys []string, signalType types.SignalType) ([]*types.Field, error) {
	return nil, nil
}

func (s *clickhouseFieldMetadataStore) GetAllFieldsMetadata(signalType types.SignalType) ([]*types.Field, error) {
	return nil, nil
}

func (s *clickhouseFieldMetadataStore) GetSampleValues(fieldKey *types.FieldKey, signalType types.SignalType, limit int) ([]any, error) {
	return nil, nil
}

func (s *clickhouseFieldMetadataStore) GetRelatedValues(fieldKey *types.FieldKey, value any, signalType types.SignalType, limit int) ([]any, error) {
	return nil, nil
}
