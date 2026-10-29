package graalsystems

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider_InternalValidate(t *testing.T) {
	require.NoError(t, Provider(DefaultProviderConfig())().InternalValidate())
}

func TestAuthSettings_ResolveMode(t *testing.T) {
	mode, err := authSettings{applicationID: "app", applicationSecret: "secret"}.resolveMode()
	require.NoError(t, err)
	assert.Equal(t, authModeApplication, mode)

	mode, err = authSettings{username: "ada", password: "pwd"}.resolveMode()
	require.NoError(t, err)
	assert.Equal(t, authModeCredentials, mode)

	_, err = authSettings{}.resolveMode()
	assert.ErrorContains(t, err, "no credentials")

	_, err = authSettings{mode: authModeApplication, applicationID: "app"}.resolveMode()
	assert.ErrorContains(t, err, "application_secret")
}

// fakePlatform serves the realm lookup, the token endpoint and a few API calls.
func fakePlatform(t *testing.T, patches *[]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/tenants/acme":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"realm":"acme-realm"}`)
		case r.URL.Path == "/realms/acme-realm/protocol/openid-connect/token":
			require.NoError(t, r.ParseForm())
			assert.Equal(t, "client_credentials", r.Form.Get("grant_type"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"access_token":"t0k3n","token_type":"Bearer","expires_in":300}`)
		case r.Method == http.MethodPatch && strings.HasPrefix(r.URL.Path, "/api/v1/jobs/"):
			assert.Equal(t, "Bearer t0k3n", r.Header.Get("Authorization"))
			assert.Equal(t, "acme", r.Header.Get("X-Tenant"))
			assert.Equal(t, "application/json-patch+json;charset=UTF-8", r.Header.Get("Content-Type"))
			body, _ := io.ReadAll(r.Body)
			*patches = append(*patches, string(body))
			_, _ = io.WriteString(w, `{}`)
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestBuildApi_ApplicationAuthAndJSONPatch(t *testing.T) {
	var patches []string
	srv := fakePlatform(t, &patches)
	defer srv.Close()

	ctx := context.Background()
	client, err := buildApi(ctx, authSettings{
		apiURL:            srv.URL + "/api/v1",
		authURL:           srv.URL,
		tenant:            "acme",
		applicationID:     "ci-bot",
		applicationSecret: "secret",
	}, "test")
	require.NoError(t, err)

	meta := &Meta{apiClient: client, tenant: "acme"}
	ops := []patchOperation{
		{Op: "replace", Path: "/timeout_seconds", Value: 600},
		{Op: "replace", Path: "/labels", Value: map[string]interface{}{"team": "data"}},
		{Op: "remove", Path: "/schedule"},
	}
	require.NoError(t, applyPatch(ctx, meta, "/jobs/j1", ops))
	require.Len(t, patches, 1)

	var sent []map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(patches[0]), &sent))
	assert.Equal(t, []map[string]interface{}{
		{"op": "replace", "path": "/timeout_seconds", "value": float64(600)},
		{"op": "replace", "path": "/labels", "value": map[string]interface{}{"team": "data"}},
		{"op": "remove", "path": "/schedule"},
	}, sent)

	err = applyPatch(ctx, meta, "/projects/unknown", ops)
	assert.ErrorContains(t, err, "HTTP 404")
}

func TestBuildApi_WrongCredentials(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error":"invalid_client"}`)
	}))
	defer srv.Close()

	_, err := buildApi(context.Background(), authSettings{
		apiURL: srv.URL + "/api/v1", authURL: srv.URL, tenant: "acme",
		applicationID: "ci-bot", applicationSecret: "wrong",
	}, "test")
	assert.ErrorContains(t, err, "authentication against")
}

func TestReplaceChanged(t *testing.T) {
	r := resourceGraalSystemsProject()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"name": "Sales", "description": "Analytics"})
	ops := replaceChanged(d, "name", "description")
	assert.Equal(t, []patchOperation{
		{Op: "replace", Path: "/name", Value: "Sales"},
		{Op: "replace", Path: "/description", Value: "Analytics"},
	}, ops)
}

func TestMaxRetriesAcceptsZero(t *testing.T) {
	validate := resourceGraalSystemsJob().Schema["max_retries"].ValidateFunc
	_, errs := validate(0, "max_retries")
	assert.Empty(t, errs)
	_, errs = validate(-1, "max_retries")
	assert.NotEmpty(t, errs)
}

func TestDefineOptions(t *testing.T) {
	block := map[string]interface{}{
		"type":          "python",
		"file":          "first_job.py",
		"instance_type": "Standard_Development_D1_v1",
		"docker_image":  "",
		"module":        "",
		"lines":         []interface{}{},
		"env":           map[string]interface{}{"STAGE": "prod"},
	}
	assert.Nil(t, validateOptions(block))
	assert.Equal(t, map[string]interface{}{
		"type":          "python",
		"file":          "first_job.py",
		"instance_type": "Standard_Development_D1_v1",
		"env":           map[string]string{"STAGE": "prod"},
	}, defineOptions(block))

	assert.NotNil(t, validateOptions(map[string]interface{}{"type": "python"}))
	assert.NotNil(t, validateOptions(map[string]interface{}{"type": "bash", "lines": []interface{}{}}))
	assert.Nil(t, validateOptions(map[string]interface{}{"type": "spark", "main_class_name": "org.acme.Job"}))
}
