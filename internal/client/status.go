// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// Status represents SABnzbd status information.
type Status struct {
	Version       string         `json:"version"`
	Paused        bool           `json:"paused"`
	Speedlimit    string         `json:"speedlimit"`
	SpeedlimitAbs string         `json:"speedlimit_abs"`
	HaveWarnings  string         `json:"have_warnings"`
	Diskspace1    string         `json:"diskspace1"`
	Diskspace2    string         `json:"diskspace2"`
	Servers       []ServerStatus `json:"servers"`
}

// ServerStatus represents the status of a news server.
type ServerStatus struct {
	ServerName      string `json:"servername"`
	ServerActive    bool   `json:"serveractive"`
	ServerError     string `json:"servererror"`
	ServerPriority  int    `json:"serverpriority"`
	ServerActiveConn int   `json:"serveractiveconn"`
	ServerTotalConn  int   `json:"servertotalconn"`
}

// FullStatus represents the full status response.
type FullStatus struct {
	Status Status `json:"status"`
}

// GetStatus retrieves the current SABnzbd status.
func (c *Client) GetStatus(ctx context.Context) (*Status, error) {
	params := url.Values{}
	params.Set("mode", "status")

	var resp FullStatus
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return nil, fmt.Errorf("getting status: %w", err)
	}

	return &resp.Status, nil
}

// GetVersion retrieves the SABnzbd version.
func (c *Client) GetVersion(ctx context.Context) (string, error) {
	params := url.Values{}
	params.Set("mode", "version")

	var resp struct {
		Version string `json:"version"`
	}
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return "", fmt.Errorf("getting version: %w", err)
	}

	return resp.Version, nil
}

// GetScripts retrieves all available scripts.
func (c *Client) GetScripts(ctx context.Context) ([]string, error) {
	params := url.Values{}
	params.Set("mode", "get_scripts")

	var resp struct {
		Scripts []string `json:"scripts"`
	}
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return nil, fmt.Errorf("getting scripts: %w", err)
	}

	return resp.Scripts, nil
}
