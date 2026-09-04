package tools

import (
	"testing"
)

func TestPriority3DeleteTools(t *testing.T) {
	tools := []string{
		"certificate_delete",
		"replication_delete",
		"replication_endpoint_delete",
	}

	for _, tool := range tools {
		t.Run(tool+"/requires confirmed", func(t *testing.T) {
			cs, cleanup := connectTestServer(t, &mockTruenasClient{})
			defer cleanup()

			res := callTool(t, cs, tool, map[string]any{"id": 1, "confirmed": false})
			assertError(t, res, "confirmed must be true")
		})

		t.Run(tool+"/invalid id", func(t *testing.T) {
			cs, cleanup := connectTestServer(t, &mockTruenasClient{})
			defer cleanup()

			res := callTool(t, cs, tool, map[string]any{"id": 0, "confirmed": true})
			assertError(t, res, "id must be a positive integer")
		})

		t.Run(tool+"/succeeds", func(t *testing.T) {
			cs, cleanup := connectTestServer(t, &mockTruenasClient{})
			defer cleanup()

			res := callTool(t, cs, tool, map[string]any{"id": 1, "confirmed": true})
			assertResultJSON(t, res)
		})
	}
}
