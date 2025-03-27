package telemetrymetadata

import (
	"context"
	"fmt"

	"github.com/SigNoz/signoz/pkg/errors"
	"github.com/SigNoz/signoz/pkg/telemetrystore"
	"github.com/SigNoz/signoz/pkg/types"
)

var (
	ErrFailedToGetTracesKeys   = errors.Newf(errors.TypeInternal, errors.CodeInternal, "failed to get traces keys")
	ErrFailedToGetTblStatement = errors.Newf(errors.TypeInternal, errors.CodeInternal, "failed to get tbl statement")
)

type telemetryMetaStore struct {
	telemetrystore         telemetrystore.TelemetryStore
	tracesDBName           string
	tracesMetadataTblName  string
	indexV3TblName         string
	metricsDBName          string
	metricsMetadataTblName string
	logsDBName             string
	logsMetadataTblName    string
	relatedMetadataDBName  string
	relatedMetadataTblName string
}

func NewTelemetryMetaStore(
	telemetrystore telemetrystore.TelemetryStore,
	tracesDBName string,
	tracesMetadataTblName string,
	indexV3TblName string,
	metricsDBName string,
	metricsMetadataTblName string,
	logsDBName string,
	logsMetadataTblName string,
	relatedMetadataDBName string,
	relatedMetadataTblName string,
) types.Metadata {
	return &telemetryMetaStore{
		telemetrystore:         telemetrystore,
		tracesDBName:           tracesDBName,
		tracesMetadataTblName:  tracesMetadataTblName,
		indexV3TblName:         indexV3TblName,
		metricsDBName:          metricsDBName,
		metricsMetadataTblName: metricsMetadataTblName,
		logsDBName:             logsDBName,
		logsMetadataTblName:    logsMetadataTblName,
		relatedMetadataDBName:  relatedMetadataDBName,
		relatedMetadataTblName: relatedMetadataTblName,
	}
}

func extractFieldKeysFromTblStatement(statement string) ([]types.TelemetryFieldKey, error) {
	return nil, nil
}

func (t *telemetryMetaStore) tblStatementToFieldKeys(ctx context.Context) ([]types.TelemetryFieldKey, error) {
	query := fmt.Sprintf("SHOW CREATE TABLE %s.%s", t.tracesDBName, t.indexV3TblName)
	statements := []types.ShowCreateTableStatement{}
	err := t.telemetrystore.ClickHouseDB().Select(ctx, &statements, query)
	if err != nil {
		return nil, ErrFailedToGetTblStatement
	}

	return nil, nil
}

func (t *telemetryMetaStore) getTracesKeys(ctx context.Context, fieldKeySelector types.FieldKeySelector) ([]types.TelemetryFieldKey, error) {
	query := fmt.Sprintf(`
		SELECT
			tag_key, tag_type, tag_data_type, max(priority) as priority
		FROM (
			SELECT
				tag_key, tag_type, tag_data_type,
				CASE
					WHEN tag_type = 'spanfield' THEN 1
					WHEN tag_type = 'resource' THEN 2
					WHEN tag_type = 'scope' THEN 3
					WHEN tag_type = 'tag' THEN 4
					ELSE 5
				END as priority
			FROM %s.%s
			WHERE tag_key $1 $2
		)
		GROUP BY tag_key, tag_type, tag_data_type
		ORDER BY priority
		LIMIT $3`,
		t.tracesDBName,
		t.tracesMetadataTblName,
	)

	args := []any{}

	// if the selector type is exact, we need to match the exact key
	if fieldKeySelector.SelectorType == types.FieldKeySelectorTypeExact {
		args = append(args, "=")
		args = append(args, fieldKeySelector.Name)
	} else {
		args = append(args, "ILIKE")
		args = append(args, fmt.Sprintf("%%%s%%", fieldKeySelector.Name))
	}

	args = append(args, fieldKeySelector.Limit)

	rows, err := t.telemetrystore.ClickHouseDB().Query(ctx, query, args...)
	if err != nil {
		return nil, ErrFailedToGetTracesKeys
	}

	keys := []types.TelemetryFieldKey{}
	for rows.Next() {
		var name, typ, dataType string
		err = rows.Scan(&name, &typ, &dataType)
		if err != nil {
			return nil, ErrFailedToGetTracesKeys
		}

		keys = append(keys, types.TelemetryFieldKey{
			Name:          name,
			FieldContext:  types.FieldContextFromString(typ),
			FieldDataType: types.FieldDataTypeFromString(dataType),
		})
	}

	if rows.Err() != nil {
		return nil, ErrFailedToGetTracesKeys
	}

	return keys, nil
}

func (t *telemetryMetaStore) getMetricsKeys(ctx context.Context, fieldKeySelector types.FieldKeySelector) ([]types.TelemetryFieldKey, error) {
	return nil, nil
}

func (t *telemetryMetaStore) getLogsKeys(ctx context.Context, fieldKeySelector types.FieldKeySelector) ([]types.TelemetryFieldKey, error) {
	return nil, nil
}

func (t *telemetryMetaStore) GetKeys(ctx context.Context, fieldKeySelector types.FieldKeySelector) (map[string][]types.TelemetryFieldKey, error) {
	return nil, nil
}

func (t *telemetryMetaStore) GetKey(ctx context.Context, fieldKeySelector types.FieldKeySelector) ([]types.TelemetryFieldKey, error) {
	return nil, nil
}

func (t *telemetryMetaStore) GetRelatedValues(ctx context.Context, fieldKeySelector types.FieldKeySelector, existingSelections []types.ExistingFieldSelection) (any, error) {
	return nil, nil
}

func (t *telemetryMetaStore) GetAllValues(ctx context.Context, fieldKeySelector types.FieldKeySelector) (any, error) {
	return nil, nil
}
