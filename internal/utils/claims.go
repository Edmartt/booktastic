package utils

import (
	"context"
	"fmt"
	"slices"
	"strings"
)

type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c *CustomClaims) Validate(ctx context.Context) error {
	if c.Scope == "" {
		return nil
	}

	if strings.TrimSpace(c.Scope) != c.Scope {
		return fmt.Errorf("scope claim has invalid whitespace")
	}

	if strings.Contains(c.Scope, "  ") {
		return fmt.Errorf("scope claim contains double spaces")
	}

	return nil
}

func (c *CustomClaims) HasScope(expectedScope string) bool {
	if c.Scope == "" {
		return false
	}

	return slices.Contains(strings.Fields(c.Scope), expectedScope)
}
