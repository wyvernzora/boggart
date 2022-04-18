package model

import (
	"github.com/aws/aws-sdk-go/service/sts"
	"time"
)

type Response struct {
	Success     bool                 `yaml:"Success"`
	Error       *string              `yaml:"Error"`
	Credentials *ResponseCredentials `yaml:"Credentials"`
}

type ResponseCredentials struct {
	AccessKeyId     string `yaml:"AccessKeyId"`
	SecretAccessKey string `yaml:"SecretAccessKey"`
	SessionToken    string `yaml:"SessionToken"`
	ExpiresAt       string `yaml:"ExpiresAt"`
}

func NewSuccessResponse(credentials *sts.Credentials) *Response {
	return &Response{
		Success: true,
		Credentials: &ResponseCredentials{
			AccessKeyId:     *credentials.AccessKeyId,
			SecretAccessKey: *credentials.SecretAccessKey,
			SessionToken:    *credentials.SessionToken,
			ExpiresAt:       credentials.Expiration.Format(time.RFC3339),
		},
	}
}

func NewErrorResponse(err error) *Response {
	errorString := err.Error()
	return &Response{
		Success: false,
		Error:   &errorString,
	}
}
