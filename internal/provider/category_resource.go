// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/apatheticriku/terraform-provider-sabnzbd/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CategoryResource{}
var _ resource.ResourceWithImportState = &CategoryResource{}

func NewCategoryResource() resource.Resource {
	return &CategoryResource{}
}

// CategoryResource defines the resource implementation.
type CategoryResource struct {
	client *client.Client
}

// CategoryResourceModel describes the resource data model.
type CategoryResourceModel struct {
	Name     types.String `tfsdk:"name"`
	Dir      types.String `tfsdk:"dir"`
	Script   types.String `tfsdk:"script"`
	Priority types.Int64  `tfsdk:"priority"`
	PP       types.String `tfsdk:"pp"`
	Order    types.Int64  `tfsdk:"order"`
}

func (r *CategoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_category"
}

func (r *CategoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a download category in SABnzbd. Categories allow you to organize downloads " +
			"and apply different settings (scripts, priorities, post-processing) to different types of content.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The unique name of the category. Use `*` for the default category.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dir": schema.StringAttribute{
				MarkdownDescription: "The relative or absolute path for completed downloads in this category. " +
					"Leave empty to use the default complete folder.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"script": schema.StringAttribute{
				MarkdownDescription: "The post-processing script to run for downloads in this category. " +
					"Use `None` for no script, or `Default` to use the global default.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("None"),
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The default priority for downloads in this category. " +
					"Values: -100=Default, -2=Paused, -1=Low, 0=Normal, 1=High, 2=Force.",
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(-100),
			},
			"pp": schema.StringAttribute{
				MarkdownDescription: "Post-processing options. Values: ``=Default, `0`=None, `1`=+Repair, " +
					"`2`=+Repair/Unpack, `3`=+Repair/Unpack/Delete.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"order": schema.Int64Attribute{
				MarkdownDescription: "The display order of this category in the UI.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
		},
	}
}

func (r *CategoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CategoryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &client.CategoryInput{
		Name:     data.Name.ValueString(),
		Dir:      data.Dir.ValueString(),
		Script:   data.Script.ValueString(),
		Priority: int(data.Priority.ValueInt64()),
		PP:       data.PP.ValueString(),
		Order:    int(data.Order.ValueInt64()),
	}

	if err := r.client.SetCategory(ctx, input); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create category, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created category resource", map[string]interface{}{"name": data.Name.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CategoryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	category, err := r.client.GetCategory(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read category, got error: %s", err))
		return
	}

	data.Dir = types.StringValue(category.Dir)
	data.Script = types.StringValue(category.Script)
	data.Priority = types.Int64Value(int64(category.Priority))
	data.PP = types.StringValue(category.PP)
	data.Order = types.Int64Value(int64(category.Order))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CategoryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &client.CategoryInput{
		Name:     data.Name.ValueString(),
		Dir:      data.Dir.ValueString(),
		Script:   data.Script.ValueString(),
		Priority: int(data.Priority.ValueInt64()),
		PP:       data.PP.ValueString(),
		Order:    int(data.Order.ValueInt64()),
	}

	if err := r.client.SetCategory(ctx, input); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update category, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated category resource", map[string]interface{}{"name": data.Name.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CategoryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteCategory(ctx, data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete category, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted category resource", map[string]interface{}{"name": data.Name.ValueString()})
}

func (r *CategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
