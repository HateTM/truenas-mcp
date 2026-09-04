package truenas

import (
	"context"
	"errors"
	"fmt"
)

// Certificate represents a TLS certificate managed by TrueNAS.
type Certificate struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Certificate     string   `json:"certificate,omitempty"`
	PrivateKey      string   `json:"privatekey,omitempty"`
	CSR             string   `json:"CSR,omitempty"`
	CommonName      string   `json:"common,omitempty"`
	SAN             []string `json:"san,omitempty"`
	Country         string   `json:"country,omitempty"`
	State           string   `json:"state,omitempty"`
	City            string   `json:"city,omitempty"`
	Organization    string   `json:"organization,omitempty"`
	Email           string   `json:"email,omitempty"`
	KeyLength       int      `json:"key_length,omitempty"`
	ECCurve         string   `json:"ec_curve,omitempty"`
	DigestAlgorithm string   `json:"digest_algorithm,omitempty"`
	Lifetime        int      `json:"lifetime,omitempty"`
}

// CreateCertificateParams holds fields for creating or updating a certificate.
// CreateType selects the creation mode, e.g. CERTIFICATE_CREATE_INTERNAL,
// CERTIFICATE_CREATE_IMPORTED, CERTIFICATE_CREATE_CSR, or CERTIFICATE_CREATE_ACME.
type CreateCertificateParams struct {
	Name            string   `json:"name"`
	CreateType      string   `json:"create_type"`
	Certificate     string   `json:"certificate,omitempty"`
	PrivateKey      string   `json:"privatekey,omitempty"`
	CSR             string   `json:"CSR,omitempty"`
	CommonName      string   `json:"common,omitempty"`
	SAN             []string `json:"san,omitempty"`
	Country         string   `json:"country,omitempty"`
	State           string   `json:"state,omitempty"`
	City            string   `json:"city,omitempty"`
	Organization    string   `json:"organization,omitempty"`
	Email           string   `json:"email,omitempty"`
	KeyLength       int      `json:"key_length,omitempty"`
	ECCurve         string   `json:"ec_curve,omitempty"`
	DigestAlgorithm string   `json:"digest_algorithm,omitempty"`
	Lifetime        int      `json:"lifetime,omitempty"`
}

// ListCertificates lists configured TLS certificates.
func (c *Client) ListCertificates(ctx context.Context, opts ...ListOptions) ([]Certificate, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	var result []Certificate
	if err := c.call(ctx, "certificate.query", buildQueryParams(nil, o), &result); err != nil {
		return nil, fmt.Errorf("listing certificates: %w", err)
	}
	return result, nil
}

// GetCertificate returns a single certificate by ID.
func (c *Client) GetCertificate(ctx context.Context, id int) (*Certificate, error) {
	var result Certificate
	if err := c.call(ctx, "certificate.get_instance", []any{id}, &result); err != nil {
		return nil, fmt.Errorf("getting certificate %d: %w", id, err)
	}
	return &result, nil
}

// CreateCertificate creates a new certificate (internal, imported, CSR, or ACME).
func (c *Client) CreateCertificate(ctx context.Context, p *CreateCertificateParams) (*Certificate, error) {
	if p == nil {
		return nil, errors.New("create certificate: params required")
	}
	var result Certificate
	if err := c.call(ctx, "certificate.create", []any{p}, &result); err != nil {
		return nil, fmt.Errorf("creating certificate %q: %w", p.Name, err)
	}
	return &result, nil
}

// UpdateCertificate updates an existing certificate's editable fields (e.g. name).
func (c *Client) UpdateCertificate(ctx context.Context, id int, p *CreateCertificateParams) (*Certificate, error) {
	if p == nil {
		return nil, errors.New("update certificate: params required")
	}
	var result Certificate
	if err := c.call(ctx, "certificate.update", []any{id, p}, &result); err != nil {
		return nil, fmt.Errorf("updating certificate %d: %w", id, err)
	}
	return &result, nil
}

// DeleteCertificate deletes a certificate.
func (c *Client) DeleteCertificate(ctx context.Context, id int) error {
	if err := c.call(ctx, "certificate.delete", []any{id}, nil); err != nil {
		return fmt.Errorf("deleting certificate %d: %w", id, err)
	}
	return nil
}

// CertificateACMEServerChoices returns the ACME directory servers available for ACME-issued certificates.
func (c *Client) CertificateACMEServerChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "certificate.acme_server_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing certificate acme server choices: %w", err)
	}
	return choices, nil
}

// CertificateCountryChoices returns the valid country codes for a certificate's subject.
func (c *Client) CertificateCountryChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "certificate.country_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing certificate country choices: %w", err)
	}
	return choices, nil
}

// CertificateECCurveChoices returns the elliptic curves available for EC key generation.
func (c *Client) CertificateECCurveChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "certificate.ec_curve_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing certificate ec curve choices: %w", err)
	}
	return choices, nil
}

// CertificateExtendedKeyUsageChoices returns the valid extended key usage (EKU) values for a certificate.
func (c *Client) CertificateExtendedKeyUsageChoices(ctx context.Context) (map[string]string, error) {
	var choices map[string]string
	if err := c.call(ctx, "certificate.extended_key_usage_choices", nil, &choices); err != nil {
		return nil, fmt.Errorf("listing certificate extended key usage choices: %w", err)
	}
	return choices, nil
}
