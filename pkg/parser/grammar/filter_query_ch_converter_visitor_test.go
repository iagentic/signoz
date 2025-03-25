package parser

import (
	"fmt"
	"testing"
)

func TestConvertToClickHouse(t *testing.T) {
	queries := []string{
		`service.name="redis" http.status_code = 200 error`,
		`http.status_code = 200`,
		`paid_user=true`,
	}
	expected := []string{
		`service.name="redis" AND http.status_code = 200 AND error`,
		`http.status_code = 200`,
		`paid_user=true`,
	}
	for idx, query := range queries {
		chQuery, err := PrepareWhereClause(query)
		if err != nil {
			t.Errorf("Error converting query to ClickHouse: %v", err)
		}
		fmt.Println("query[", idx, "]", query)
		fmt.Println("expected[", idx, "]", expected[idx])
		fmt.Println("chQuery[", idx, "]", chQuery)
	}
}
