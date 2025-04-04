package telemetryspans

import (
	schema "github.com/SigNoz/signoz-otel-collector/cmd/signozschemamigrator/schema_migrator"
)

var mainColumns = map[string]schema.Column{
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
