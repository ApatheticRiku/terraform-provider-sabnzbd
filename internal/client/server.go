// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ServerInput represents the input for creating/updating a server.
type ServerInput struct {
	Name        string
	Host        string
	Port        int
	Username    string
	Password    string
	Connections int
	SSL         bool
	SSLVerify   int
	SSLCiphers  string
	Enable      bool
	Optional    bool
	Retention   int
	Timeout     int
	Priority    int
	Required    bool
	Notes       string
}

// SetServer creates or updates a news server configuration.
func (c *Client) SetServer(ctx context.Context, input *ServerInput) error {
	params := url.Values{}
	params.Set("mode", "set_config")
	params.Set("section", "servers")
	params.Set("name", input.Name)
	params.Set("host", input.Host)
	params.Set("port", strconv.Itoa(input.Port))
	params.Set("username", input.Username)
	params.Set("password", input.Password)
	params.Set("connections", strconv.Itoa(input.Connections))
	params.Set("ssl", boolToInt(input.SSL))
	params.Set("ssl_verify", strconv.Itoa(input.SSLVerify))
	params.Set("ssl_ciphers", input.SSLCiphers)
	params.Set("enable", boolToInt(input.Enable))
	params.Set("optional", boolToInt(input.Optional))
	params.Set("retention", strconv.Itoa(input.Retention))
	params.Set("timeout", strconv.Itoa(input.Timeout))
	params.Set("priority", strconv.Itoa(input.Priority))
	params.Set("required", boolToInt(input.Required))
	params.Set("notes", input.Notes)

	var resp map[string]interface{}
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return fmt.Errorf("setting server config: %w", err)
	}

	return nil
}

// GetServer retrieves a specific server configuration by name.
func (c *Client) GetServer(ctx context.Context, name string) (*Server, error) {
	config, err := c.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	for _, server := range config.Servers {
		if server.Name == name {
			return &server, nil
		}
	}

	return nil, fmt.Errorf("server %q not found", name)
}

// DeleteServer removes a server configuration.
func (c *Client) DeleteServer(ctx context.Context, name string) error {
	params := url.Values{}
	params.Set("mode", "del_config")
	params.Set("section", "servers")
	params.Set("keyword", name)

	var resp map[string]interface{}
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return fmt.Errorf("deleting server config: %w", err)
	}

	return nil
}

func boolToInt(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
