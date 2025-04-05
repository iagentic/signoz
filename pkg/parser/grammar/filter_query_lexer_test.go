package parser

import (
	"testing"
)

var queries = []string{
	// `service.name="redis"`,
	// `service.name=redis`,
	// `http.status_code = 200`,
	// `paid_user=true`,
	// `service.name!="frontend"`,
	// `http.status_code!=200`,
	// `paid_user!=true`,
	// `duration_ms < 200`,
	// `bytes_sent < 1000000`,
	// `duration_ms <= 200`,
	// `bytes_sent <= 1000000`,
	// `duration_ms > 200`,
	// `bytes_sent > 1000000`,
	// `duration_ms >= 200`,
	// `bytes_sent >= 1000000`,
	// `email like "%signoz.io%"`,
	// `user.name LIKE "%srikanth%"`,
	// `email not like "%signoz.io%"`,
	// `user.name NOT LIKE "%srikanth%"`,
	// `email ilike "%signoz.io%"`,
	// `user.name ILIKE "%srikanth%"`,
	// `http.status_code between 200 and 300`,
	// `duration_ms between 100 and 200`,
	// `bytes_sent between 100000 and 200000`,
	// `http.status_code not between 200 and 300`,
	// `duration_ms not between 100 and 200`,
	// `bytes_sent not between 100000 and 200000`,
	// `service.name in ("redis", "mysql")`,
	// `service.name not in ("redis", "mysql") http.status_code = 200`,
	// `body contains "error" k8s.pod.name="pod-123"`,
	// `body contains "error" k8s.pod.name="pod-123" AND k8s.namespace="default"`,
	// `body not contains "error" k8s.pod.name="pod-123" AND k8s.namespace="default"`,
	// `body like '%gmail.com%' http.status_code = 200 service.name!="frontend"`,
	// `body not like '%gmail.com%' http.status_code = 200 service.name!="frontend"`,
	// `"Waiting for the response"`,
	// `"Waiting for the response" http.status_code = 200`,
	// `"Waiting for the response" OR "no error" http.status_code = 200 AND http.verb="GET"`,
	// `"Waiting for the response" OR "no error" http.status_code = 200 AND http.verb="GET"`,
	// `has(payload.user_ids, 123)`,
	// `has(proto.user_objects[].name, "srikanth")`,
	// `hasAny(payload.user_ids, [123, 456])`,
	// `hasAll(proto.user_objects[].name, ["srikanth", "karthik"])`,
	// `hasNone(payload.user_ids, [123, 456])`,
	// `hasAny(payload.user_ids, [123, 456])`,
	// `hasAny(proto.user_objects[].name, ["srikanth", "karthik"])`,
	// `service.name="redis" AND http.status_code = 200`,
	// `service.name="redis" OR http.status_code = 200`,
	// `NOT service.name="redis"`,
	// `NOT http.status_code = 200`,
	// `NOT (service.name="redis" AND http.status_code = 200)`,
	// `(service.name="redis" AND http.status_code = 200) OR (service.name="mysql" AND http.status_code = 200)`,
	// `(service.name="redis"
	// AND http.status_code = 200)
	// OR (service.name="mysql"
	// AND http.status_code = 200)`,
	// `(service.name=redis
	// AND http.status_code = 200)
	// OR (service.name=mysql
	// AND http.status_code = 200)`,
	`error message http.status_code = 200`,
	// `service.name IN ("redis", "mysql")`,
	// `service.name NOT IN ("redis", "mysql") http.status_code = 200`,
	// `body contains "error" k8s.pod.name="pod-123"`,
	// `syntax error: ']' came as a complete surprise to me 343$#$@#@#^$ 454%$^^&% #@#@ aksa :a;aa;s 121`,
	// `[traceId=fd1dc36b345b4296a5dfce66449b71d4; spanId=ab72922b5290c943; jobId=audience:5800-sync:42794; method=null]`,
	"boghy>=nil",
}

func TestLexer(t *testing.T) {

}
