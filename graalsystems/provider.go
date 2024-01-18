package graalsystems

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"os"

	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var debug = os.Getenv("GS_DEBUG") != ""

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
				"username": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The username (for credentials auth mode).",
				},
				"password": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The password (for credentials auth mode).",
				},
				"application_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The application id (for application auth mode).",
				},
				"application_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The application secret (for application auth mode).",
				},
				"tenant": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The tenant ID.",
				},
				"api_url": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The API URL to use.",
				},
				"auth_url": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The Auth URL to use.",
				},
				"auth_mode": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The Auth mode to use.",
				},
			},

			ResourcesMap: map[string]*schema.Resource{
				"graalsystems_project":   resourceGraalSystemsProject(),
				"graalsystems_identity":  resourceGraalSystemsIdentity(),
				"graalsystems_job":       resourceGraalSystemsJob(),
				"graalsystems_user":      resourceGraalSystemsUser(),
				"graalsystems_group":     resourceGraalSystemsGroup(),
				"graalsystems_workspace": resourceGraalSystemsWorkspace(),
			},

			DataSourcesMap: map[string]*schema.Resource{
				"graalsystems_project":   dataSourceGraalSystemsProject(),
				"graalsystems_identity":  dataSourceGraalSystemsIdentity(),
				"graalsystems_job":       dataSourceGraalSystemsJob(),
				"graalsystems_user":      dataSourceGraalSystemsUser(),
				"graalsystems_group":     dataSourceGraalSystemsGroup(),
				"graalsystems_workspace": dataSourceGraalSystemsWorkspace(),
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
	httpClient       *http.Client
}

// providerConfigure creates the Meta object containing the SDK client.
func buildMeta(ctx context.Context, config *metaConfig) (*Meta, error) {
	tenant := config.providerSchema.Get("tenant").(string)
	apiUrl := config.providerSchema.Get("api_url").(string)
	authUrl := config.providerSchema.Get("auth_url").(string)
	username := config.providerSchema.Get("username").(string)
	password := config.providerSchema.Get("password").(string)
	applicationId := config.providerSchema.Get("application_id").(string)
	applicationSecret := config.providerSchema.Get("application_secret").(string)
	authMode := config.providerSchema.Get("auth_mode").(string)
	terraformVersion := config.terraformVersion

	apiClient, err := buildApi(ctx, apiUrl, authUrl, terraformVersion, tenant, username, password, applicationId, applicationSecret, authMode)
	if err != nil {
		return nil, err
	}

	return &Meta{
		apiClient: apiClient,
		tenant:    tenant,
	}, nil
}

func buildApi(ctx context.Context, apiUrl string, authUrl string, terraformVersion string, tenant string, username string, password string, appId string, appSecret string, authMode string) (*sdk.APIClient, error) {
	////
	// Create GraalSystems SDK client
	////
	servers := sdk.ServerConfigurations{}
	servers = append(servers, sdk.ServerConfiguration{
		URL: apiUrl,
	})

	authUrl, err := findRealm(ctx, terraformVersion, servers, tenant, authUrl)
	if err != nil {
		return nil, err
	}

	var client *http.Client

	if authMode == "" || authMode == "credentials" {
		cfg := oauth2.Config{
			ClientID: "graal-ui",
			Endpoint: oauth2.Endpoint{
				TokenURL: authUrl,
			},
		}
		client, _ = buildOAuth2ClientCredentials(ctx, cfg, username, password)
	} else if authMode == "application" {
		cfg := clientcredentials.Config{
			ClientID:     appId,
			ClientSecret: appSecret,
			TokenURL:     authUrl,
		}
		client, _ = buildOAuth2ClientApplication(ctx, cfg)
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid auth mode: %s", authMode))
	}

	configuration := sdk.Configuration{
		UserAgent:  fmt.Sprintf("terraform-provider/%s terraform/%s", version, terraformVersion),
		Debug:      debug,
		HTTPClient: client,
		Servers:    servers,
	}

	apiClient := sdk.NewAPIClient(&configuration)
	return apiClient, nil
}

func findRealm(ctx context.Context, terraformVersion string, servers sdk.ServerConfigurations, tenant string, authUrl string) (string, error) {
	tmpClient := &http.Client{}
	tmpConfiguration := sdk.Configuration{
		UserAgent:  fmt.Sprintf("terraform-provider/%s terraform/%s", version, terraformVersion),
		Debug:      debug,
		HTTPClient: tmpClient,
		Servers:    servers,
	}
	tmpApiClient := sdk.NewAPIClient(&tmpConfiguration)
	t, _, err := tmpApiClient.TenantAPI.FindRealmByTenantId(ctx, tenant).Execute()
	if err != nil {
		return "", errors.WithStack(err)
	}
	authUrl = authUrl + "/realms/" + *t.Realm + "/protocol/openid-connect/token"
	return authUrl, nil
}

func buildOAuth2ClientCredentials(ctx context.Context, cfg oauth2.Config, username string, password string) (*http.Client, error) {
	token, err := cfg.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !token.Valid() {
		return nil, errors.New(fmt.Sprintf("Token invalid. Got: %#v", token))
	}

	var client *http.Client
	if debug {
		trace := buildClientTrace()
		client = cfg.Client(httptrace.WithClientTrace(ctx, trace), token)
	} else {
		client = cfg.Client(ctx, token)
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, client)
	return client, nil
}

func buildOAuth2ClientApplication(ctx context.Context, cfg clientcredentials.Config) (*http.Client, error) {
	var client *http.Client
	if debug {
		trace := buildClientTrace()
		client = cfg.Client(httptrace.WithClientTrace(ctx, trace))
	} else {
		client = cfg.Client(ctx)
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, client)
	_, err := cfg.Token(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client, nil
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
