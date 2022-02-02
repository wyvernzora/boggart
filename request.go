package main

import (
	"encoding/json"
	"fmt"
	"github.com/gliderlabs/ssh"
)

type Request struct {
	RoleArn     string `json:"roleArn"`
	SessionName string `json:"sessionName"`
}

func ParseRequest(s ssh.Session) (*Request, error) {
	var request Request
	err := json.Unmarshal([]byte(s.RawCommand()), &request)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request JSON: %w", err)
	}

	return &request, nil
}
