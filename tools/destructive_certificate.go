package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerCertificateDestructiveTools adds the opt-in destructive certificate.delete tool
// to the server. Only registered when TRUENAS_ALLOW_DESTRUCTIVE=true.
func registerCertificateDestructiveTools(s *mcp.Server, client truenasClient) {
	destructiveHint := true

	type certificateIDInput struct {
		ID        int  `json:"id"        jsonschema:"Numeric certificate ID"`
		Confirmed bool `json:"confirmed" jsonschema:"Must be set to true to confirm deletion"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "certificate_delete",
		Description: "Permanently delete a TLS certificate. Fails if it is still in use by a service. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p certificateIDInput) (*mcp.CallToolResult, any, error) {
		if !p.Confirmed {
			return errorResult(errors.New("certificate_delete: confirmed must be true to proceed with deletion"))
		}
		if p.ID <= 0 {
			return errorResult(errors.New("certificate_delete: id must be a positive integer"))
		}
		if err := client.DeleteCertificate(ctx, p.ID); err != nil {
			return errorResult(fmt.Errorf("certificate_delete: %w", err))
		}
		return jsonResult(map[string]any{"deleted": true, "id": p.ID})
	})
}
