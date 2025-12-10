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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FoldersResource{}
var _ resource.ResourceWithImportState = &FoldersResource{}

func NewFoldersResource() resource.Resource {
	return &FoldersResource{}
}

// FoldersResource defines the resource implementation.
type FoldersResource struct {
	client *client.Client
}

// FoldersResourceModel describes the resource data model.
type FoldersResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	DownloadDir         types.String `tfsdk:"download_dir"`
	DownloadFree        types.String `tfsdk:"download_free"`
	CompleteDir         types.String `tfsdk:"complete_dir"`
	CompleteFree        types.String `tfsdk:"complete_free"`
	AutoResume          types.Bool   `tfsdk:"auto_resume"`
	Permissions         types.String `tfsdk:"permissions"`
	WatchedDir          types.String `tfsdk:"watched_dir"`
	WatchedDirScanSpeed types.Int64  `tfsdk:"watched_dir_scan_speed"`
	ScriptsDir          types.String `tfsdk:"scripts_dir"`
	EmailTemplatesDir   types.String `tfsdk:"email_templates_dir"`
	PasswordFile        types.String `tfsdk:"password_file"`
	NzbBackupDir        types.String `tfsdk:"nzb_backup_dir"`
	AdminDir            types.String `tfsdk:"admin_dir"`
	BackupDir           types.String `tfsdk:"backup_dir"`
	LogDir              types.String `tfsdk:"log_dir"`
}

func (r *FoldersResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folders"
}

func (r *FoldersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages folder configuration in SABnzbd. This resource configures paths for downloads, " +
			"scripts, watched folders, and other directory settings. Note: This is a singleton resource - only one instance should exist.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier (always 'folders').",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"download_dir": schema.StringAttribute{
				MarkdownDescription: "Temporary download folder where files are stored during download. " +
					"Can be relative to base folder (e.g., 'Incomplete') or absolute path.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"download_free": schema.StringAttribute{
				MarkdownDescription: "Minimum free space for temporary download folder (e.g., '10G', '500M'). " +
					"SABnzbd pauses when free space falls below this value.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"complete_dir": schema.StringAttribute{
				MarkdownDescription: "Completed download folder for finished downloads. " +
					"This is the default location unless overridden by categories.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"complete_free": schema.StringAttribute{
				MarkdownDescription: "Minimum free space for completed download folder (e.g., '10G', '500M'). " +
					"SABnzbd pauses when free space falls below this value.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"auto_resume": schema.BoolAttribute{
				MarkdownDescription: "Automatically resume downloading when minimum free space becomes available again.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"permissions": schema.StringAttribute{
				MarkdownDescription: "Permissions for completed downloads in octal notation (e.g., '755', '777'). " +
					"Only applies to macOS and Linux.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"watched_dir": schema.StringAttribute{
				MarkdownDescription: "Folder periodically scanned for new NZB files. " +
					"Supports category sub-folders and filename prefixes for automatic categorization.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"watched_dir_scan_speed": schema.Int64Attribute{
				MarkdownDescription: "Seconds between filesystem scans of watched folder. Set to 0 to disable automatic scans.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(5),
			},
			"scripts_dir": schema.StringAttribute{
				MarkdownDescription: "Folder where user scripts (post-processing and pre-queue) are stored.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"email_templates_dir": schema.StringAttribute{
				MarkdownDescription: "Folder for custom email templates.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"password_file": schema.StringAttribute{
				MarkdownDescription: "Path to text file containing known passwords (one per line) for passworded RAR files.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"nzb_backup_dir": schema.StringAttribute{
				MarkdownDescription: "Folder where NZB files are backed up after processing.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"admin_dir": schema.StringAttribute{
				MarkdownDescription: "Folder for SABnzbd administrative files.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("admin"),
			},
			"backup_dir": schema.StringAttribute{
				MarkdownDescription: "Folder for SABnzbd configuration backups.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("backup"),
			},
			"log_dir": schema.StringAttribute{
				MarkdownDescription: "Folder for SABnzbd log files.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("logs"),
			},
		},
	}
}

func (r *FoldersResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FoldersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FoldersResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &client.FoldersInput{
		DownloadDir:         data.DownloadDir.ValueString(),
		DownloadFree:        data.DownloadFree.ValueString(),
		CompleteDir:         data.CompleteDir.ValueString(),
		CompleteFree:        data.CompleteFree.ValueString(),
		AutoResume:          data.AutoResume.ValueBool(),
		Permissions:         data.Permissions.ValueString(),
		WatchedDir:          data.WatchedDir.ValueString(),
		WatchedDirScanSpeed: int(data.WatchedDirScanSpeed.ValueInt64()),
		ScriptsDir:          data.ScriptsDir.ValueString(),
		EmailTemplatesDir:   data.EmailTemplatesDir.ValueString(),
		PasswordFile:        data.PasswordFile.ValueString(),
		NzbBackupDir:        data.NzbBackupDir.ValueString(),
		AdminDir:            data.AdminDir.ValueString(),
		BackupDir:           data.BackupDir.ValueString(),
		LogDir:              data.LogDir.ValueString(),
	}

	if err := r.client.SetFolders(ctx, input); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create folders configuration, got error: %s", err))
		return
	}

	data.ID = types.StringValue("folders")
	tflog.Trace(ctx, "created folders resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FoldersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FoldersResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	folders, err := r.client.GetFolders(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read folders configuration, got error: %s", err))
		return
	}

	data.ID = types.StringValue("folders")
	data.DownloadDir = types.StringValue(folders.DownloadDir)
	data.DownloadFree = types.StringValue(folders.DownloadFree)
	data.CompleteDir = types.StringValue(folders.CompleteDir)
	data.CompleteFree = types.StringValue(folders.CompleteFree)
	data.AutoResume = types.BoolValue(folders.AutoResume == 1)
	data.Permissions = types.StringValue(folders.Permissions)
	data.WatchedDir = types.StringValue(folders.WatchedDir)
	data.WatchedDirScanSpeed = types.Int64Value(int64(folders.WatchedDirScanSpeed))
	data.ScriptsDir = types.StringValue(folders.ScriptsDir)
	data.EmailTemplatesDir = types.StringValue(folders.EmailTemplatesDir)
	data.PasswordFile = types.StringValue(folders.PasswordFile)
	data.NzbBackupDir = types.StringValue(folders.NzbBackupDir)
	data.AdminDir = types.StringValue(folders.AdminDir)
	data.BackupDir = types.StringValue(folders.BackupDir)
	data.LogDir = types.StringValue(folders.LogDir)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FoldersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FoldersResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &client.FoldersInput{
		DownloadDir:         data.DownloadDir.ValueString(),
		DownloadFree:        data.DownloadFree.ValueString(),
		CompleteDir:         data.CompleteDir.ValueString(),
		CompleteFree:        data.CompleteFree.ValueString(),
		AutoResume:          data.AutoResume.ValueBool(),
		Permissions:         data.Permissions.ValueString(),
		WatchedDir:          data.WatchedDir.ValueString(),
		WatchedDirScanSpeed: int(data.WatchedDirScanSpeed.ValueInt64()),
		ScriptsDir:          data.ScriptsDir.ValueString(),
		EmailTemplatesDir:   data.EmailTemplatesDir.ValueString(),
		PasswordFile:        data.PasswordFile.ValueString(),
		NzbBackupDir:        data.NzbBackupDir.ValueString(),
		AdminDir:            data.AdminDir.ValueString(),
		BackupDir:           data.BackupDir.ValueString(),
		LogDir:              data.LogDir.ValueString(),
	}

	if err := r.client.SetFolders(ctx, input); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update folders configuration, got error: %s", err))
		return
	}

	data.ID = types.StringValue("folders")
	tflog.Trace(ctx, "updated folders resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FoldersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Folders configuration cannot be deleted, only reset to defaults
	// We'll just remove it from state
	tflog.Trace(ctx, "deleted folders resource from state")
}

func (r *FoldersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
