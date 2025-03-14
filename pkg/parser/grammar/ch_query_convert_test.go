package parser

import (
	"testing"

	"go.signoz.io/signoz/pkg/types"
)

func TestPrepareWhereClause(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple equality",
			input:    `status = "active"`,
			expected: `status = "active"`,
			wantErr:  false,
		},
		{
			name:     "Numeric comparison",
			input:    `count > 100`,
			expected: `count > 100`,
			wantErr:  false,
		},
		{
			name:     "AND condition",
			input:    `status = "active" AND user_id = 123`,
			expected: `status = "active" AND user_id = 123`,
			wantErr:  false,
		},
		{
			name:     "OR condition",
			input:    `status = "active" OR status = "pending"`,
			expected: `status = "active" OR status = "pending"`,
			wantErr:  false,
		},
		{
			name:     "Mixed AND/OR with parentheses",
			input:    `(status = "active" OR status = "pending") AND created_at > "2023-01-01"`,
			expected: `(status = "active" OR status = "pending") AND created_at > "2023-01-01"`,
			wantErr:  false,
		},
		{
			name:     "NOT condition",
			input:    `NOT is_deleted = 1`,
			expected: `NOT (is_deleted = 1)`,
			wantErr:  false,
		},
		{
			name:     "NOT with comparison",
			input:    `NOT (status = "inactive")`,
			expected: `NOT (status = "inactive")`,
			wantErr:  false,
		},
		{
			name:     "IN clause with parentheses",
			input:    `status IN ("active", "pending")`,
			expected: `status IN ("active", "pending")`,
			wantErr:  false,
		},
		{
			name:     "IN clause with brackets",
			input:    `status IN ["active", "pending"]`,
			expected: `status IN ("active", "pending")`,
			wantErr:  false,
		},
		{
			name:     "NOT IN clause",
			input:    `status NOT IN ("deleted", "archived")`,
			expected: `status NOT IN ("deleted", "archived")`,
			wantErr:  false,
		},
		{
			name:     "BETWEEN clause",
			input:    `created_at BETWEEN "2023-01-01" AND "2023-12-31"`,
			expected: `created_at BETWEEN "2023-01-01" AND "2023-12-31"`,
			wantErr:  false,
		},
		{
			name:     "NOT BETWEEN clause",
			input:    `created_at NOT BETWEEN "2023-01-01" AND "2023-12-31"`,
			expected: `(created_at < "2023-01-01" OR created_at > "2023-12-31")`,
			wantErr:  false,
		},
		{
			name:     "LIKE clause",
			input:    `name LIKE "%smith%"`,
			expected: `name LIKE "%smith%"`,
			wantErr:  false,
		},
		{
			name:     "NOT LIKE clause",
			input:    `name NOT LIKE "%test%"`,
			expected: `name NOT LIKE "%test%"`,
			wantErr:  false,
		},
		{
			name:     "ILIKE clause",
			input:    `name ILIKE "%smith%"`,
			expected: `name ILIKE "%smith%"`,
			wantErr:  false,
		},
		{
			name:     "Function call - has",
			input:    `has(tags, "important")`,
			expected: `has(tags, "important")`,
			wantErr:  false,
		},
		{
			name:     "Function call - hasAny",
			input:    `hasAny(categories, ["billing", "payment"])`,
			expected: `hasAny(categories, ["billing", "payment"])`,
			wantErr:  false,
		},
		{
			name:     "Function call - hasAll",
			input:    `hasAll(tags, ["urgent", "important"])`,
			expected: `hasAll(tags, ["urgent", "important"])`,
			wantErr:  false,
		},
		{
			name:     "Function call - hasNone",
			input:    `hasNone(tags, ["deleted", "archived"])`,
			expected: `not hasAny(tags, ["deleted", "archived"])`,
			wantErr:  false,
		},
		{
			name:     "Full text search",
			input:    `"error connecting to database"`,
			expected: `lower(body) LIKE '%%error connecting to database%%' OR has(mapValues(attributes_string, '%%error connecting to database%%'))`,
			wantErr:  false,
		},
		{
			name:     "Complex query with multiple conditions",
			input:    `(status = "active" OR status = "pending") AND created_at > "2023-01-01" AND has(tags, "important")`,
			expected: `(status = "active" OR status = "pending") AND created_at > "2023-01-01" AND has(tags, "important")`,
			wantErr:  false,
		},
		{
			name:     "Implicit AND condition",
			input:    `status = "active" user_id = 123`,
			expected: `(status = "active") AND (user_id = 123)`,
			wantErr:  false,
		},
		{
			// For key EXISTS conditions, they need to be specified in the grammar
			name:     "EXISTS condition",
			input:    `tags EXISTS`,
			expected: `has(tags)`,
			wantErr:  false,
		},
		{
			name:     "NOT EXISTS condition",
			input:    `tags NOT EXISTS`,
			expected: `not has(tags)`,
			wantErr:  false,
		},
		{
			name:     "REGEXP condition",
			input:    `name REGEXP "^user_[0-9]+"`,
			expected: `match(name, "^user_[0-9]+")`,
			wantErr:  false,
		},
		{
			name:     "NOT REGEXP condition",
			input:    `name NOT REGEXP "^test_"`,
			expected: `not match(name, "^test_")`,
			wantErr:  false,
		},
		{
			name:     "CONTAINS condition",
			input:    `description CONTAINS "error"`,
			expected: `position(description, "error") > 0`,
			wantErr:  false,
		},
		{
			name:     "NOT CONTAINS condition",
			input:    `description NOT CONTAINS "success"`,
			expected: `position(description, "success") = 0`,
			wantErr:  false,
		},
		{
			name:    "Invalid syntax",
			input:   `status = `,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrepareWhereClause(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("PrepareWhereClause() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.expected {
				t.Errorf("PrepareWhereClause() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestQBFieldFromKeyText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *types.QBField
	}{
		{
			name:     "service.name",
			input:    "service.name",
			expected: &types.QBField{Key: "service.name", Context: types.QBFieldContextResource, DataType: types.QBFieldDataTypeString},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QBFieldFromKeyText(tt.input)
			if got != tt.expected {
				t.Errorf("QBFieldFromKeyText() = %v, want %v", got, tt.expected)
			}
		})
	}
}
