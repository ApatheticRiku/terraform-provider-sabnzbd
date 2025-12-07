// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/apatheticriku/terraform-provider-sabnzbd/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure SabnzbdProvider satisfies various provider interfaces.
var _ provider.Provider = &SabnzbdProvider{}

// SabnzbdProvider defines the provider implementation.
type SabnzbdProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SabnzbdProviderModel describes the provider data model.
type SabnzbdProviderModel struct {
	URL    types.String `tfsdk:"url"`
	APIKey types.String `tfsdk:"api_key"`
}

func (p *SabnzbdProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sabnzbd"
	resp.Version = p.version
}

func (p *SabnzbdProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The SABnzbd provider allows you to manage SABnzbd configuration as infrastructure. " +
			"It supports managing servers, categories, and other settings through the SABnzbd API.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the SABnzbd instance (e.g., `http://localhost:8080`). " +
					"Can also be set via the `SABNZBD_URL` environment variable.",
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for authenticating with SABnzbd. " +
					"Can also be set via the `SABNZBD_API_KEY` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *SabnzbdProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SabnzbdProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Default to environment variables if not set in configuration.
	url := os.Getenv("SABNZBD_URL")
	apiKey := os.Getenv("SABNZBD_API_KEY")

	if !data.URL.IsNull() {
		url = data.URL.ValueString()
	}

	if !data.APIKey.IsNull() {
		apiKey = data.APIKey.ValueString()
	}

	// Validate required configuration.
	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing SABnzbd URL",
			"The provider cannot create the SABnzbd API client as there is a missing or empty value for the SABnzbd URL. "+
				"Set the url value in the configuration or use the SABNZBD_URL environment variable.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing SABnzbd API Key",
			"The provider cannot create the SABnzbd API client as there is a missing or empty value for the SABnzbd API key. "+
				"Set the api_key value in the configuration or use the SABNZBD_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the SABnzbd client.
	sabnzbdClient := client.NewClient(url, apiKey)

	resp.DataSourceData = sabnzbdClient
	resp.ResourceData = sabnzbdClient
}

func (p *SabnzbdProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewServerResource,
		NewCategoryResource,
	}
}

func (p *SabnzbdProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewConfigDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SabnzbdProvider{
			version: version,
		}
	}
}
