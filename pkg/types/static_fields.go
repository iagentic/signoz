package types

var TraceStaticFieldsTraces = map[string]FieldKey{
	"timestamp": {
		Key:      "timestamp",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"span_id": {
		Key:      "span_id",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"trace_state": {
		Key:      "trace_state",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"parent_span_id": {
		Key:      "parent_span_id",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"flags": {
		Key:      "flags",
		DataType: FieldDataTypeNumber,
		Context:  FieldContextPrimary,
	},
	"name": {
		Key:      "name",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"kind": {
		Key:      "kind",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"kind_string": {
		Key:      "kind_string",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"duration_nano": {
		Key:      "duration_nano",
		DataType: FieldDataTypeNumber,
		Context:  FieldContextPrimary,
	},
	"status_code": {
		Key:      "status_code",
		DataType: FieldDataTypeNumber,
		Context:  FieldContextPrimary,
	},
	"status_message": {
		Key:      "status_message",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"status_code_string": {
		Key:      "status_code_string",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},

	// new support for composite attributes
	"response_status_code": {
		Key:      "response_status_code",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"external_http_url": {
		Key:      "external_http_url",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"http_url": {
		Key:      "http_url",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"external_http_method": {
		Key:      "external_http_method",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"http_method": {
		Key:      "http_method",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"http_host": {
		Key:      "http_host",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"db_name": {
		Key:      "db_name",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"db_operation": {
		Key:      "db_operation",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
	"has_error": {
		Key:      "has_error",
		DataType: FieldDataTypeBoolean,
		Context:  FieldContextPrimary,
	},
	"is_remote": {
		Key:      "is_remote",
		DataType: FieldDataTypeString,
		Context:  FieldContextPrimary,
	},
}
