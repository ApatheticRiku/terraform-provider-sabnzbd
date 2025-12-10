// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// FoldersInput represents the input for updating folder configuration.
type FoldersInput struct {
	DownloadDir         string
	DownloadFree        string
	CompleteDir         string
	CompleteFree        string
	AutoResume          bool
	Permissions         string
	WatchedDir          string
	WatchedDirScanSpeed int
	ScriptsDir          string
	EmailTemplatesDir   string
	PasswordFile        string
	NzbBackupDir        string
	AdminDir            string
	BackupDir           string
	LogDir              string
}

// Folders represents the folder configuration from SABnzbd.
type Folders struct {
	DownloadDir         string `json:"download_dir"`
	DownloadFree        string `json:"download_free"`
	CompleteDir         string `json:"complete_dir"`
	CompleteFree        string `json:"complete_free"`
	AutoResume          int    `json:"auto_resume"`
	Permissions         string `json:"permissions"`
	WatchedDir          string `json:"dirscan_dir"`
	WatchedDirScanSpeed int    `json:"dirscan_speed"`
	ScriptsDir          string `json:"script_dir"`
	EmailTemplatesDir   string `json:"email_dir"`
	PasswordFile        string `json:"password_file"`
	NzbBackupDir        string `json:"nzb_backup_dir"`
	AdminDir            string `json:"admin_dir"`
	BackupDir           string `json:"backup_dir"`
	LogDir              string `json:"log_dir"`
}

// SetFolders updates the folder configuration.
func (c *Client) SetFolders(ctx context.Context, input *FoldersInput) error {
	params := url.Values{}
	params.Set("mode", "set_config")
	params.Set("section", "misc")

	if input.DownloadDir != "" {
		params.Set("download_dir", input.DownloadDir)
	}
	if input.DownloadFree != "" {
		params.Set("download_free", input.DownloadFree)
	}
	if input.CompleteDir != "" {
		params.Set("complete_dir", input.CompleteDir)
	}
	if input.CompleteFree != "" {
		params.Set("complete_free", input.CompleteFree)
	}
	params.Set("auto_resume", boolToInt(input.AutoResume))
	if input.Permissions != "" {
		params.Set("permissions", input.Permissions)
	}
	if input.WatchedDir != "" {
		params.Set("dirscan_dir", input.WatchedDir)
	}
	if input.WatchedDirScanSpeed >= 0 {
		params.Set("dirscan_speed", fmt.Sprintf("%d", input.WatchedDirScanSpeed))
	}
	if input.ScriptsDir != "" {
		params.Set("script_dir", input.ScriptsDir)
	}
	if input.EmailTemplatesDir != "" {
		params.Set("email_dir", input.EmailTemplatesDir)
	}
	if input.PasswordFile != "" {
		params.Set("password_file", input.PasswordFile)
	}
	if input.NzbBackupDir != "" {
		params.Set("nzb_backup_dir", input.NzbBackupDir)
	}
	if input.AdminDir != "" {
		params.Set("admin_dir", input.AdminDir)
	}
	if input.BackupDir != "" {
		params.Set("backup_dir", input.BackupDir)
	}
	if input.LogDir != "" {
		params.Set("log_dir", input.LogDir)
	}

	var resp map[string]interface{}
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return fmt.Errorf("setting folders config: %w", err)
	}

	return nil
}

// GetFolders retrieves the folder configuration.
func (c *Client) GetFolders(ctx context.Context) (*Folders, error) {
	config, err := c.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	folders := &Folders{}
	misc := config.Misc

	if v, ok := misc["download_dir"].(string); ok {
		folders.DownloadDir = v
	}
	if v, ok := misc["download_free"].(string); ok {
		folders.DownloadFree = v
	}
	if v, ok := misc["complete_dir"].(string); ok {
		folders.CompleteDir = v
	}
	if v, ok := misc["complete_free"].(string); ok {
		folders.CompleteFree = v
	}
	if v, ok := misc["auto_resume"].(float64); ok {
		folders.AutoResume = int(v)
	}
	if v, ok := misc["permissions"].(string); ok {
		folders.Permissions = v
	}
	if v, ok := misc["dirscan_dir"].(string); ok {
		folders.WatchedDir = v
	}
	if v, ok := misc["dirscan_speed"].(float64); ok {
		folders.WatchedDirScanSpeed = int(v)
	}
	if v, ok := misc["script_dir"].(string); ok {
		folders.ScriptsDir = v
	}
	if v, ok := misc["email_dir"].(string); ok {
		folders.EmailTemplatesDir = v
	}
	if v, ok := misc["password_file"].(string); ok {
		folders.PasswordFile = v
	}
	if v, ok := misc["nzb_backup_dir"].(string); ok {
		folders.NzbBackupDir = v
	}
	if v, ok := misc["admin_dir"].(string); ok {
		folders.AdminDir = v
	}
	if v, ok := misc["backup_dir"].(string); ok {
		folders.BackupDir = v
	}
	if v, ok := misc["log_dir"].(string); ok {
		folders.LogDir = v
	}

	return folders, nil
}
