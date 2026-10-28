package graalsystems

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strings"

	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var debug = os.Getenv("GRAAL_DEBUG") != "" || os.Getenv("GS_DEBUG") != ""

// ProviderConfig config can be used to provide additional config when creating provider.
type ProviderConfig struct {
	// Meta can be used to override Meta that will be used by the provider.
	// This is useful for tests.
	Meta *Meta
}

// DefaultProviderConfig return default ProviderConfig struct
func DefaultProviderConfig() *ProviderConfig {
	return &ProviderConfig{}
}

// Provider returns a terraform.ResourceProvider.
func Provider(config *ProviderConfig) plugin.ProviderFunc {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_url": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"GRAAL_API_URL", "GS_API_URL"}, nil),
					Description: "URL of the graal REST API, including `/api/v1`. Can be set with `GRAAL_API_URL`.",
				},
				"auth_url": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"GRAAL_KEYCLOAK_URL", "GS_AUTH_URL"}, nil),
					Description: "Base URL of the OpenID Connect server (Keycloak), without `/realms/...`. Can be set with `GRAAL_KEYCLOAK_URL`.",
				},
				"tenant": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"GRAAL_TENANT", "GS_TENANT"}, nil),
					Description: "The tenant, sent as the `X-Tenant` header. Can be set with `GRAAL_TENANT`.",
				},
				"realm": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("GRAAL_REALM", nil),
					Description: "The authentication realm. Discovered from the tenant when omitted. Can be set with `GRAAL_REALM`.",
				},
				"auth_mode": {
					Type:         schema.TypeString,
					Optional:     true,
					DefaultFunc:  schema.MultiEnvDefaultFunc([]string{"GRAAL_AUTH_MODE", "GS_AUTH_MODE"}, nil),
					ValidateFunc: validation.StringInSlice([]string{authModeApplication, authModeCredentials}, false),
					Description:  "`application` (client credentials, recommended for automation) or `credentials` (username and password). Inferred from the attributes that are set when omitted.",
				},
				"application_id": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"GRAAL_CLIENT_ID", "GS_APPLICATION_ID"}, nil),
					Description: "The application (client) id, for the `application` auth mode. Can be set with `GRAAL_CLIENT_ID`.",
				},
				"application_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"GRAAL_CLIENT_SECRET", "GS_APPLICATION_SECRET"}, nil),
					Description: "The application secret, for the `application` auth mode. Can be set with `GRAAL_CLIENT_SECRET`.",
				},
				"username": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"GRAAL_USERNAME", "GS_USERNAME"}, nil),
					Description: "The username, for the `credentials` auth mode. Can be set with `GRAAL_USERNAME`.",
				},
				"password": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{"GRAAL_PASSWORD", "GS_PASSWORD"}, nil),
					Description: "The password, for the `credentials` auth mode. Can be set with `GRAAL_PASSWORD`.",
				},
			},

			ResourcesMap: map[string]*schema.Resource{
				"graalsystems_project":   resourceGraalSystemsProject(),
				"graalsystems_identity":  resourceGraalSystemsIdentity(),
				"graalsystems_job":       resourceGraalSystemsJob(),
				"graalsystems_user":      resourceGraalSystemsUser(),
				"graalsystems_group":     resourceGraalSystemsGroup(),
				"graalsystems_workspace": resourceGraalSystemsWorkspace(),
				"graalsystems_workflow":  resourceGraalSystemsWorkflow(),
			},

			DataSourcesMap: map[string]*schema.Resource{
				"graalsystems_project":   dataSourceGraalSystemsProject(),
				"graalsystems_identity":  dataSourceGraalSystemsIdentity(),
				"graalsystems_job":       dataSourceGraalSystemsJob(),
				"graalsystems_user":      dataSourceGraalSystemsUser(),
				"graalsystems_group":     dataSourceGraalSystemsGroup(),
				"graalsystems_workspace": dataSourceGraalSystemsWorkspace(),
				"graalsystems_workflow":  dataSourceGraalSystemsWorkflow(),
			},
		}

		p.ConfigureContextFunc = func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
			terraformVersion := p.TerraformVersion

			// If we provide meta in config use it. This is useful for tests
			if config.Meta != nil {
				return config.Meta, nil
			}

			meta, err := buildMeta(ctx, &metaConfig{
				providerSchema:   data,
				terraformVersion: terraformVersion,
			})
			if err != nil {
				return nil, diag.FromErr(err)
			}
			return meta, nil
		}

		return p
	}
}

// Meta contains config and SDK clients used by resources.
//
// This meta value is passed into all resources.
type Meta struct {
	// apiClient is the GraalSystems SDK client.
	apiClient *sdk.APIClient
	tenant    string
}

type metaConfig struct {
	providerSchema   *schema.ResourceData
	terraformVersion string
}

// authSettings holds the authentication attributes of the provider block.
type authSettings struct {
	apiURL, authURL, tenant, realm, mode string
	username, password                   string
	applicationID, applicationSecret     string
}

const (
	authModeApplication = "application"
	authModeCredentials = "credentials"
)

// resolveMode returns the auth mode to use, inferring it when not set.
func (s authSettings) resolveMode() (string, error) {
	mode := s.mode
	if mode == "" {
		switch {
		case s.applicationID != "" || s.applicationSecret != "":
			mode = authModeApplication
		case s.username != "" || s.password != "":
			mode = authModeCredentials
		default:
			return "", errors.New("no credentials: set application_id and application_secret (or GRAAL_CLIENT_ID and GRAAL_CLIENT_SECRET), or username and password")
		}
	}
	switch mode {
	case authModeApplication:
		if s.applicationID == "" || s.applicationSecret == "" {
			return "", errors.New("the application auth mode needs application_id and application_secret")
		}
	case authModeCredentials:
		if s.username == "" || s.password == "" {
			return "", errors.New("the credentials auth mode needs username and password")
		}
	default:
		return "", fmt.Errorf("invalid auth mode: %s", mode)
	}
	return mode, nil
}

// buildMeta creates the Meta object containing the SDK client.
func buildMeta(ctx context.Context, config *metaConfig) (*Meta, error) {
	d := config.providerSchema
	settings := authSettings{
		apiURL:            d.Get("api_url").(string),
		authURL:           d.Get("auth_url").(string),
		tenant:            d.Get("tenant").(string),
		realm:             d.Get("realm").(string),
		mode:              d.Get("auth_mode").(string),
		username:          d.Get("username").(string),
		password:          d.Get("password").(string),
		applicationID:     d.Get("application_id").(string),
		applicationSecret: d.Get("application_secret").(string),
	}

	apiClient, err := buildApi(ctx, settings, config.terraformVersion)
	if err != nil {
		return nil, err
	}

	return &Meta{
		apiClient: apiClient,
		tenant:    settings.tenant,
	}, nil
}

func buildApi(ctx context.Context, s authSettings, terraformVersion string) (*sdk.APIClient, error) {
	for name, value := range map[string]string{"api_url": s.apiURL, "auth_url": s.authURL, "tenant": s.tenant} {
		if value == "" {
			return nil, fmt.Errorf("the provider attribute %s is required", name)
		}
	}
	mode, err := s.resolveMode()
	if err != nil {
		return nil, err
	}

	servers := sdk.ServerConfigurations{{URL: strings.TrimRight(s.apiURL, "/")}}
	userAgent := fmt.Sprintf("terraform-provider-graalsystems/%s terraform/%s", version, terraformVersion)

	realm := s.realm
	if realm == "" {
		realm = findRealm(ctx, userAgent, servers, s.tenant)
	}
	tokenURL := strings.TrimRight(s.authURL, "/") + "/realms/" + url.PathEscape(realm) + "/protocol/openid-connect/token"

	var client *http.Client
	if mode == authModeCredentials {
		cfg := oauth2.Config{
			ClientID: "graal-ui",
			Endpoint: oauth2.Endpoint{TokenURL: tokenURL},
		}
		client, err = buildOAuth2ClientCredentials(ctx, cfg, s.username, s.password)
	} else {
		cfg := clientcredentials.Config{
			ClientID:     s.applicationID,
			ClientSecret: s.applicationSecret,
			TokenURL:     tokenURL,
		}
		client, err = buildOAuth2ClientApplication(ctx, cfg)
	}
	if err != nil {
		return nil, fmt.Errorf("authentication against %s failed: %w", tokenURL, err)
	}

	configuration := sdk.Configuration{
		UserAgent:  userAgent,
		Debug:      debug,
		HTTPClient: client,
		Servers:    servers,
	}

	return sdk.NewAPIClient(&configuration), nil
}

// findRealm returns the authentication realm of a tenant. The lookup does not
// need authentication; when it fails, the realm is assumed to be the tenant.
func findRealm(ctx context.Context, userAgent string, servers sdk.ServerConfigurations, tenant string) string {
	tmpApiClient := sdk.NewAPIClient(&sdk.Configuration{
		UserAgent:  userAgent,
		Debug:      debug,
		HTTPClient: &http.Client{},
		Servers:    servers,
	})
	t, _, err := tmpApiClient.TenantAPI.FindRealmByTenantId(ctx, tenant).Execute()
	if err != nil || t == nil || t.Realm == nil || *t.Realm == "" {
		return tenant
	}
	return *t.Realm
}

// clientContext returns the context bound to an OAuth2 HTTP client. Tokens are
// refreshed long after the provider is configured, so it must not be the
// (short-lived) configuration context.
func clientContext() context.Context {
	if debug {
		return httptrace.WithClientTrace(context.Background(), buildClientTrace())
	}
	return context.Background()
}

func buildOAuth2ClientCredentials(ctx context.Context, cfg oauth2.Config, username string, password string) (*http.Client, error) {
	token, err := cfg.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, err
	}
	if !token.Valid() {
		return nil, errors.New("the identity provider returned an invalid token")
	}
	return cfg.Client(clientContext(), token), nil
}

func buildOAuth2ClientApplication(ctx context.Context, cfg clientcredentials.Config) (*http.Client, error) {
	// Fail early on wrong credentials rather than on the first API call.
	if _, err := cfg.Token(ctx); err != nil {
		return nil, err
	}
	return cfg.Client(clientContext()), nil
}

func buildClientTrace() *httptrace.ClientTrace {
	trace := &httptrace.ClientTrace{
		GetConn:      func(hostPort string) { fmt.Println("starting to create conn", hostPort) },
		DNSStart:     func(info httptrace.DNSStartInfo) { fmt.Println("starting to look up dns", info) },
		DNSDone:      func(info httptrace.DNSDoneInfo) { fmt.Println("done looking up dns", info) },
		ConnectStart: func(network, addr string) { fmt.Println("starting tcp connection", network, addr) },
		ConnectDone:  func(network, addr string, err error) { fmt.Println("tcp connection created", network, addr, err) },
		GotConn:      func(info httptrace.GotConnInfo) { fmt.Println("connection established", info) },
	}
	return trace
}
