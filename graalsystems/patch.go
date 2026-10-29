package graalsystems

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// patchOperation is one JSON Patch (RFC 6902) operation.
type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// replaceChanged returns an operation that sets each attribute that changed.
// Attribute names must match the API field names.
//
// The operation is "add", not "replace": the API omits empty fields from some
// resources (a user without description, for example), and "replace" on a
// missing member fails with INVALID_PATCH. "add" sets the member whether it
// exists or not (RFC 6902, section 4.1).
func replaceChanged(d *schema.ResourceData, attributes ...string) []patchOperation {
	var ops []patchOperation
	for _, a := range attributes {
		if d.HasChange(a) {
			ops = append(ops, patchOperation{Op: "add", Path: "/" + a, Value: d.Get(a)})
		}
	}
	return ops
}

// applyPatch sends a JSON Patch document to a resource of the API (there is no
// full PUT). The SDK's Patch model only carries string values, so the request
// is made directly with the SDK's authenticated HTTP client.
//
// kind names the resource type of the vendor media type the API answers with
// (application/vnd.graal.systems.v1.<kind>+json): any other Accept header is
// refused with 406 Not Acceptable.
func applyPatch(ctx context.Context, meta *Meta, kind, path string, ops []patchOperation) error {
	if len(ops) == 0 {
		return nil
	}
	cfg := meta.apiClient.GetConfig()
	if len(cfg.Servers) == 0 {
		return fmt.Errorf("no API URL configured")
	}
	body, err := json.Marshal(ops)
	if err != nil {
		return err
	}
	url := strings.TrimRight(cfg.Servers[0].URL, "/") + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json-patch+json;charset=UTF-8")
	req.Header.Set("Accept", "application/vnd.graal.systems.v1."+kind+"+json")
	req.Header.Set("X-Tenant", meta.tenant)
	req.Header.Set("User-Agent", cfg.UserAgent)

	client := cfg.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("PATCH %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("PATCH %s: HTTP %d: %s", path, resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	return nil
}
