package boggart

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gliderlabs/ssh"
	"time"
)

type Response struct {
	Success     bool         `json:"success"`
	Credentials *Credentials `json:"credentials,omitempty"`
	Error       *string      `json:"error,omitempty"`
}

type Credentials struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	ExpiresAt       string `json:"expiresAt"`
}

func SuccessResponse(credentials *sts.Credentials) *Response {
	return &Response{
		Success: true,
		Credentials: &Credentials{
			AccessKeyId:     *credentials.AccessKeyId,
			SecretAccessKey: *credentials.SecretAccessKey,
			SessionToken:    *credentials.SessionToken,
			ExpiresAt:       credentials.Expiration.Format(time.RFC3339),
		},
	}
}

func ErrorResponse(err error) *Response {
	errorString := err.Error()
	return &Response{
		Success: false,
		Error:   &errorString,
	}
}

func (response *Response) Write(s ssh.Session) {
	data, _ := json.Marshal(response)
	_, _ = s.Write(data)
}
