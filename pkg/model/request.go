package model

import (
	"fmt"
	"net/url"
	"strings"
)

type Request struct {
	RoleArn     string
	SessionName string
	Format      string
}

func ParseRequest(in string) (*Request, error) {
	var req Request

	parts := strings.SplitN(in, "?", 2)
	req.RoleArn = parts[0]

	if len(parts) == 2 {
		query, err := url.ParseQuery(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse request: %w", err)
		}
		req.SessionName = query.Get("session")
		req.Format = query.Get("format")
	}

	return &req, nil
}
