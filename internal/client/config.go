// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// GetConfigResponse represents the full configuration response.
type GetConfigResponse struct {
	Config Config `json:"config"`
}

// Config represents SABnzbd configuration sections.
type Config struct {
	Misc       map[string]interface{} `json:"misc"`
	Servers    []Server               `json:"servers"`
	Categories []Category             `json:"categories"`
	RSS        []RSSFeed              `json:"rss"`
	Sorters    []Sorter               `json:"sorters"`
}

// Server represents a news server configuration.
type Server struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Connections int    `json:"connections"`
	SSL         int    `json:"ssl"`
	SSLVerify   int    `json:"ssl_verify"`
	SSLCiphers  string `json:"ssl_ciphers"`
	Enable      int    `json:"enable"`
	Optional    int    `json:"optional"`
	Retention   int    `json:"retention"`
	Timeout     int    `json:"timeout"`
	Priority    int    `json:"priority"`
	Required    int    `json:"required"`
	Notes       string `json:"notes"`
}

// Category represents a download category configuration.
type Category struct {
	Name     string `json:"name"`
	Dir      string `json:"dir"`
	Script   string `json:"script"`
	Priority int    `json:"priority"`
	PP       string `json:"pp"`
	Order    int    `json:"order"`
}

// RSSFeed represents an RSS feed configuration.
type RSSFeed struct {
	Name     string   `json:"name"`
	URI      []string `json:"uri"`
	Cat      string   `json:"cat"`
	PP       string   `json:"pp"`
	Script   string   `json:"script"`
	Enable   int      `json:"enable"`
	Priority int      `json:"priority"`
}

// Sorter represents a sorting rule configuration.
type Sorter struct {
	Name       string   `json:"name"`
	Order      int      `json:"order"`
	SortString string   `json:"sort_string"`
	SortCats   []string `json:"sort_cats"`
	SortType   []int    `json:"sort_type"`
	IsActive   int      `json:"is_active"`
}

// GetConfig retrieves the full SABnzbd configuration.
func (c *Client) GetConfig(ctx context.Context) (*Config, error) {
	params := url.Values{}
	params.Set("mode", "get_config")

	var resp GetConfigResponse
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return nil, fmt.Errorf("getting config: %w", err)
	}

	return &resp.Config, nil
}

// GetConfigSection retrieves a specific section of the configuration.
func (c *Client) GetConfigSection(ctx context.Context, section string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("mode", "get_config")
	params.Set("section", section)

	var resp map[string]interface{}
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return nil, fmt.Errorf("getting config section %s: %w", section, err)
	}

	return resp, nil
}

// GetConfigSectionByKeyword retrieves a specific item from a section by keyword.
func (c *Client) GetConfigSectionByKeyword(ctx context.Context, section, keyword string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("mode", "get_config")
	params.Set("section", section)
	params.Set("keyword", keyword)

	var resp map[string]interface{}
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return nil, fmt.Errorf("getting config section %s keyword %s: %w", section, keyword, err)
	}

	return resp, nil
}
