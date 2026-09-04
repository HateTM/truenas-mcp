package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/gordcurrie/truenas-mcp/internal/truenas"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerCertificateTools registers certificate.* MCP tools onto the server.
func registerCertificateTools(s *mcp.Server, client truenasClient) {
	type emptyInput struct{}

	type listCertificatesInput struct {
		Limit  int `json:"limit,omitempty"  jsonschema:"Maximum number of certificates to return; 0 means no limit"`
		Offset int `json:"offset,omitempty" jsonschema:"Number of certificates to skip; 0 means start from the beginning"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "certificate_list",
		Description: "List configured TLS certificates.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p listCertificatesInput) (*mcp.CallToolResult, any, error) {
		result, err := client.ListCertificates(ctx, truenas.ListOptions{Limit: p.Limit, Offset: p.Offset})
		if err != nil {
			return errorResult(fmt.Errorf("certificate_list: %w", err))
		}
		return jsonResult(result)
	})

	type getCertificateInput struct {
		ID int `json:"id" jsonschema:"Numeric certificate ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "certificate_get",
		Description: "Get a single TLS certificate by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p getCertificateInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("certificate_get: id must be a positive integer"))
		}
		result, err := client.GetCertificate(ctx, p.ID)
		if err != nil {
			return errorResult(fmt.Errorf("certificate_get: %w", err))
		}
		return jsonResult(result)
	})

	type certificateInput struct {
		Name            string   `json:"name"                       jsonschema:"Certificate name"`
		CreateType      string   `json:"create_type"                jsonschema:"CERTIFICATE_CREATE_INTERNAL, CERTIFICATE_CREATE_IMPORTED, CERTIFICATE_CREATE_CSR, or CERTIFICATE_CREATE_ACME"`
		Certificate     string   `json:"certificate,omitempty"      jsonschema:"PEM certificate content, required when create_type=CERTIFICATE_CREATE_IMPORTED"`
		PrivateKey      string   `json:"privatekey,omitempty"       jsonschema:"PEM private key content, required when importing a certificate or CSR"`
		CSR             string   `json:"CSR,omitempty"              jsonschema:"PEM CSR content, required when create_type=CERTIFICATE_CREATE_CSR"`
		CommonName      string   `json:"common,omitempty"           jsonschema:"Subject common name (CN), required for CERTIFICATE_CREATE_INTERNAL/ACME"`
		SAN             []string `json:"san,omitempty"              jsonschema:"Subject Alternative Names"`
		Country         string   `json:"country,omitempty"          jsonschema:"Subject country code, from certificate_country_choices"`
		State           string   `json:"state,omitempty"            jsonschema:"Subject state/province"`
		City            string   `json:"city,omitempty"             jsonschema:"Subject city"`
		Organization    string   `json:"organization,omitempty"     jsonschema:"Subject organization"`
		Email           string   `json:"email,omitempty"            jsonschema:"Subject email address"`
		KeyLength       int      `json:"key_length,omitempty"       jsonschema:"RSA key length in bits, e.g. 2048"`
		ECCurve         string   `json:"ec_curve,omitempty"         jsonschema:"EC curve name, from certificate_ec_curve_choices; use instead of key_length for EC keys"`
		DigestAlgorithm string   `json:"digest_algorithm,omitempty" jsonschema:"Signature digest algorithm, e.g. SHA256"`
		Lifetime        int      `json:"lifetime,omitempty"         jsonschema:"Certificate validity period in days, for internally-generated certificates"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "certificate_create",
		Description: "Create a new TLS certificate: internally signed, imported, from a CSR, or via ACME.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p certificateInput) (*mcp.CallToolResult, any, error) {
		if p.Name == "" || p.CreateType == "" {
			return errorResult(errors.New("certificate_create: name and create_type are required"))
		}
		result, err := client.CreateCertificate(ctx, &truenas.CreateCertificateParams{
			Name: p.Name, CreateType: p.CreateType, Certificate: p.Certificate, PrivateKey: p.PrivateKey,
			CSR: p.CSR, CommonName: p.CommonName, SAN: p.SAN, Country: p.Country, State: p.State,
			City: p.City, Organization: p.Organization, Email: p.Email, KeyLength: p.KeyLength,
			ECCurve: p.ECCurve, DigestAlgorithm: p.DigestAlgorithm, Lifetime: p.Lifetime,
		})
		if err != nil {
			return errorResult(fmt.Errorf("certificate_create: %w", err))
		}
		return jsonResult(result)
	})

	type updateCertificateInput struct {
		ID   int    `json:"id"   jsonschema:"Numeric certificate ID"`
		Name string `json:"name" jsonschema:"New certificate name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "certificate_update",
		Description: "Update an existing certificate's name.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, p updateCertificateInput) (*mcp.CallToolResult, any, error) {
		if p.ID <= 0 {
			return errorResult(errors.New("certificate_update: id must be a positive integer"))
		}
		if p.Name == "" {
			return errorResult(errors.New("certificate_update: name must not be empty"))
		}
		result, err := client.UpdateCertificate(ctx, p.ID, &truenas.CreateCertificateParams{Name: p.Name})
		if err != nil {
			return errorResult(fmt.Errorf("certificate_update: %w", err))
		}
		return jsonResult(result)
	})

	// certificate_delete is destructive and lives in destructive_certificate.go,
	// gated behind Config.AllowDestructive.

	mcp.AddTool(s, &mcp.Tool{
		Name:        "certificate_acme_server_choices",
		Description: "List the ACME directory servers available for ACME-issued certificates.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.CertificateACMEServerChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("certificate_acme_server_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "certificate_country_choices",
		Description: "List the valid country codes for a certificate's subject.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.CertificateCountryChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("certificate_country_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "certificate_ec_curve_choices",
		Description: "List the elliptic curves available for EC key generation.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.CertificateECCurveChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("certificate_ec_curve_choices: %w", err))
		}
		return jsonResult(choices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "certificate_extended_key_usage_choices",
		Description: "List the valid extended key usage (EKU) values for a certificate.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, any, error) {
		choices, err := client.CertificateExtendedKeyUsageChoices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("certificate_extended_key_usage_choices: %w", err))
		}
		return jsonResult(choices)
	})
}
