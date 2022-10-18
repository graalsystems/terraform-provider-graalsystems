package graalsystems

import (
	"context"
	"fmt"
	sdk "github.com/graalsystems/sdk/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/http/httptrace"
	"os"
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
					Description: "The username.",
				},
				"password": {
					Type:        schema.TypeString,
					Optional:    true, // To allow user to use deprecated `token`.
					Description: "The password.",
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
				"graalsystems_project":  resourceGraalSystemsProject(),
				"graalsystems_identity": resourceGraalSystemsIdentity(),
				"graalsystems_job":      resourceGraalSystemsJob(),
			},

			DataSourcesMap: map[string]*schema.Resource{
				"graalsystems_project":  dataSourceGraalSystemsProject(),
				"graalsystems_identity": dataSourceGraalSystemsIdentity(),
				"graalsystems_job":      dataSourceGraalSystemsJob(),
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
	////
	// Create GraalSystems SDK client
	////
	servers := sdk.ServerConfigurations{}
	servers = append(servers, sdk.ServerConfiguration{
		URL: config.providerSchema.Get("api_url").(string),
	})

	cfg := clientcredentials.Config{
		ClientID:     config.providerSchema.Get("username").(string),
		ClientSecret: config.providerSchema.Get("password").(string),
		TokenURL:     config.providerSchema.Get("auth_url").(string),
	}

	var client *http.Client
	if debug {
		trace := &httptrace.ClientTrace{
			GetConn:      func(hostPort string) { fmt.Println("starting to create conn ", hostPort) },
			DNSStart:     func(info httptrace.DNSStartInfo) { fmt.Println("starting to look up dns", info) },
			DNSDone:      func(info httptrace.DNSDoneInfo) { fmt.Println("done looking up dns", info) },
			ConnectStart: func(network, addr string) { fmt.Println("starting tcp connection", network, addr) },
			ConnectDone:  func(network, addr string, err error) { fmt.Println("tcp connection created", network, addr, err) },
			GotConn:      func(info httptrace.GotConnInfo) { fmt.Println("connection established", info) },
		}
		client = cfg.Client(httptrace.WithClientTrace(context.Background(), trace))
	} else {
		client = cfg.Client(context.Background())
	}

	configuration := sdk.Configuration{
		UserAgent:  fmt.Sprintf("terraform-provider/%s terraform/%s", version, config.terraformVersion),
		Debug:      debug,
		HTTPClient: client,
		Servers:    servers,
	}

	apiClient := sdk.NewAPIClient(&configuration)

	return &Meta{
		apiClient: apiClient,
		tenant:    config.providerSchema.Get("tenant").(string),
	}, nil
}
