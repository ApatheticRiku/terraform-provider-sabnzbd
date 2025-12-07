// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/apatheticriku/terraform-provider-sabnzbd/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ConfigDataSource{}

func NewConfigDataSource() datasource.DataSource {
	return &ConfigDataSource{}
}

// ConfigDataSource defines the data source implementation.
type ConfigDataSource struct {
	client *client.Client
}

// ConfigDataSourceModel describes the data source data model.
type ConfigDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Version    types.String `tfsdk:"version"`
	Categories types.List   `tfsdk:"categories"`
	Scripts    types.List   `tfsdk:"scripts"`
}

func (d *ConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (d *ConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves basic configuration information from SABnzbd, including version, " +
			"available categories, and available scripts.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier for this data source.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The version of SABnzbd.",
				Computed:            true,
			},
			"categories": schema.ListAttribute{
				MarkdownDescription: "List of available category names.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"scripts": schema.ListAttribute{
				MarkdownDescription: "List of available script names.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *ConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *ConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConfigDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get version.
	version, err := d.client.GetVersion(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read version, got error: %s", err))
		return
	}
	data.Version = types.StringValue(version)

	// Get categories.
	categories, err := d.client.GetCategories(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read categories, got error: %s", err))
		return
	}
	categoryValues := make([]types.String, len(categories))
	for i, cat := range categories {
		categoryValues[i] = types.StringValue(cat)
	}
	categoriesList, diags := types.ListValueFrom(ctx, types.StringType, categoryValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Categories = categoriesList

	// Get scripts.
	scripts, err := d.client.GetScripts(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scripts, got error: %s", err))
		return
	}
	scriptValues := make([]types.String, len(scripts))
	for i, script := range scripts {
		scriptValues[i] = types.StringValue(script)
	}
	scriptsList, diags := types.ListValueFrom(ctx, types.StringType, scriptValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Scripts = scriptsList

	data.ID = types.StringValue("sabnzbd-config")

	tflog.Trace(ctx, "read config data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
