// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/apatheticriku/terraform-provider-sabnzbd/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ServerResource{}
var _ resource.ResourceWithImportState = &ServerResource{}

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

// ServerResource defines the resource implementation.
type ServerResource struct {
	client *client.Client
}

// ServerResourceModel describes the resource data model.
type ServerResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Host        types.String `tfsdk:"host"`
	Port        types.Int64  `tfsdk:"port"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Connections types.Int64  `tfsdk:"connections"`
	SSL         types.Bool   `tfsdk:"ssl"`
	SSLVerify   types.Int64  `tfsdk:"ssl_verify"`
	SSLCiphers  types.String `tfsdk:"ssl_ciphers"`
	Enable      types.Bool   `tfsdk:"enable"`
	Optional    types.Bool   `tfsdk:"optional"`
	Retention   types.Int64  `tfsdk:"retention"`
	Timeout     types.Int64  `tfsdk:"timeout"`
	Priority    types.Int64  `tfsdk:"priority"`
	Required    types.Bool   `tfsdk:"required"`
	Notes       types.String `tfsdk:"notes"`
}

func (r *ServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *ServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a news server configuration in SABnzbd.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The unique name/identifier for this server configuration.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The hostname or IP address of the news server.",
				Required:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port number for the news server. Default is 563 for SSL, 119 for non-SSL.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(563),
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username for authentication.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password for authentication.",
				Optional:            true,
				Sensitive:           true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"connections": schema.Int64Attribute{
				MarkdownDescription: "The number of connections to use for this server.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(8),
			},
			"ssl": schema.BoolAttribute{
				MarkdownDescription: "Whether to use SSL/TLS for the connection.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"ssl_verify": schema.Int64Attribute{
				MarkdownDescription: "SSL certificate verification level: 0=None, 1=Verify CA, 2=Verify CA and hostname.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2),
			},
			"ssl_ciphers": schema.StringAttribute{
				MarkdownDescription: "Custom SSL ciphers to use (leave empty for default).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Whether this server is enabled.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"optional": schema.BoolAttribute{
				MarkdownDescription: "Whether this server is optional (used only when primary servers fail).",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"retention": schema.Int64Attribute{
				MarkdownDescription: "The retention period in days (0 for unlimited).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "Connection timeout in seconds.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Server priority (0 is highest priority).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Whether this server is required for downloads to complete.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Optional notes about this server.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
		},
	}
}

func (r *ServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &client.ServerInput{
		Name:        data.Name.ValueString(),
		Host:        data.Host.ValueString(),
		Port:        int(data.Port.ValueInt64()),
		Username:    data.Username.ValueString(),
		Password:    data.Password.ValueString(),
		Connections: int(data.Connections.ValueInt64()),
		SSL:         data.SSL.ValueBool(),
		SSLVerify:   int(data.SSLVerify.ValueInt64()),
		SSLCiphers:  data.SSLCiphers.ValueString(),
		Enable:      data.Enable.ValueBool(),
		Optional:    data.Optional.ValueBool(),
		Retention:   int(data.Retention.ValueInt64()),
		Timeout:     int(data.Timeout.ValueInt64()),
		Priority:    int(data.Priority.ValueInt64()),
		Required:    data.Required.ValueBool(),
		Notes:       data.Notes.ValueString(),
	}

	if err := r.client.SetServer(ctx, input); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create server, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created server resource", map[string]interface{}{"name": data.Name.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.GetServer(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read server, got error: %s", err))
		return
	}

	data.Host = types.StringValue(server.Host)
	data.Port = types.Int64Value(int64(server.Port))
	data.Username = types.StringValue(server.Username)
	// Note: Password is not returned by the API for security reasons.
	data.Connections = types.Int64Value(int64(server.Connections))
	data.SSL = types.BoolValue(server.SSL == 1)
	data.SSLVerify = types.Int64Value(int64(server.SSLVerify))
	data.SSLCiphers = types.StringValue(server.SSLCiphers)
	data.Enable = types.BoolValue(server.Enable == 1)
	data.Optional = types.BoolValue(server.Optional == 1)
	data.Retention = types.Int64Value(int64(server.Retention))
	data.Timeout = types.Int64Value(int64(server.Timeout))
	data.Priority = types.Int64Value(int64(server.Priority))
	data.Required = types.BoolValue(server.Required == 1)
	data.Notes = types.StringValue(server.Notes)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &client.ServerInput{
		Name:        data.Name.ValueString(),
		Host:        data.Host.ValueString(),
		Port:        int(data.Port.ValueInt64()),
		Username:    data.Username.ValueString(),
		Password:    data.Password.ValueString(),
		Connections: int(data.Connections.ValueInt64()),
		SSL:         data.SSL.ValueBool(),
		SSLVerify:   int(data.SSLVerify.ValueInt64()),
		SSLCiphers:  data.SSLCiphers.ValueString(),
		Enable:      data.Enable.ValueBool(),
		Optional:    data.Optional.ValueBool(),
		Retention:   int(data.Retention.ValueInt64()),
		Timeout:     int(data.Timeout.ValueInt64()),
		Priority:    int(data.Priority.ValueInt64()),
		Required:    data.Required.ValueBool(),
		Notes:       data.Notes.ValueString(),
	}

	if err := r.client.SetServer(ctx, input); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update server, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated server resource", map[string]interface{}{"name": data.Name.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteServer(ctx, data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete server, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted server resource", map[string]interface{}{"name": data.Name.ValueString()})
}

func (r *ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
