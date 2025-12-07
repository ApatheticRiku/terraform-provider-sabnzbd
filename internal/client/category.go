// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// CategoryInput represents the input for creating/updating a category.
type CategoryInput struct {
	Name     string
	Dir      string
	Script   string
	Priority int
	PP       string
	Order    int
}

// SetCategory creates or updates a category configuration.
func (c *Client) SetCategory(ctx context.Context, input *CategoryInput) error {
	params := url.Values{}
	params.Set("mode", "set_config")
	params.Set("section", "categories")
	params.Set("name", input.Name)
	params.Set("dir", input.Dir)
	params.Set("script", input.Script)
	params.Set("priority", strconv.Itoa(input.Priority))
	params.Set("pp", input.PP)
	params.Set("order", strconv.Itoa(input.Order))

	var resp map[string]interface{}
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return fmt.Errorf("setting category config: %w", err)
	}

	return nil
}

// GetCategory retrieves a specific category configuration by name.
func (c *Client) GetCategory(ctx context.Context, name string) (*Category, error) {
	config, err := c.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	for _, cat := range config.Categories {
		if cat.Name == name {
			return &cat, nil
		}
	}

	return nil, fmt.Errorf("category %q not found", name)
}

// DeleteCategory removes a category configuration.
func (c *Client) DeleteCategory(ctx context.Context, name string) error {
	params := url.Values{}
	params.Set("mode", "del_config")
	params.Set("section", "categories")
	params.Set("keyword", name)

	var resp map[string]interface{}
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return fmt.Errorf("deleting category config: %w", err)
	}

	return nil
}

// GetCategories retrieves all category names.
func (c *Client) GetCategories(ctx context.Context) ([]string, error) {
	params := url.Values{}
	params.Set("mode", "get_cats")

	var resp struct {
		Categories []string `json:"categories"`
	}
	if err := c.doRequest(ctx, params, &resp); err != nil {
		return nil, fmt.Errorf("getting categories: %w", err)
	}

	return resp.Categories, nil
}
