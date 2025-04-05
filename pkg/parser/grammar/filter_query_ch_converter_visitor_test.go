package parser

import (
	"fmt"
	"testing"
)

func TestConvertToClickHouse(t *testing.T) {
	queries := []string{
		// `service.name="redis" http.status_code = 200 error`,
		// `http.status_code = 200`,
		// `paid_user=true`,
		// `"[traceId=fd1dc36b345b4296a5dfce66449b71d4; spanId=ab72922b5290c943; jobId=audience:5800-sync:42794; method=null]"`,
		// "error message http.status_code = 200",
		// "waiting_for_response",
		"boghy>=nil",
	}
	expected := []string{
		// `service.name="redis" AND http.status_code = 200 AND error`,
		// `http.status_code = 200`,
		// `paid_user=true`,
		// `"[traceId=fd1dc36b345b4296a5dfce66449b71d4; spanId=ab72922b5290c943; jobId=audience:5800-sync:42794; method=null]"`,
		// `error AND http.status_code = 200`,
		// `waiting_for_response`,
		`boghy>=nil`,
	}
	for idx, query := range queries {
		chQuery, err := PrepareWhereClause(query, nil, nil)
		if err != nil {
			t.Errorf("Error converting query to ClickHouse: %v", err)
		}
		fmt.Println("query[", idx, "]", query)
		fmt.Println("expected[", idx, "]", expected[idx])
		fmt.Println("chQuery[", idx, "]", chQuery)
	}
}
