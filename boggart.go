package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gliderlabs/ssh"
)

type Boggart struct {
	stsClient *sts.STS
	config    *Config
}

func NewBoggart(configfile string) (*Boggart, error) {
	config, err := LoadConfig(configfile)
	if err != nil {
		return nil, err
	}
	LogWithoutContext("service.config_loaded", map[string]interface{}{})

	awsSession := session.Must(session.NewSession(&aws.Config{
		Region:      config.Region,
		Credentials: config.Credentials,
	}))
	stsClient := sts.New(awsSession)
	LogWithoutContext("service.sts_ready", map[string]interface{}{})

	boggart := &Boggart{
		stsClient: stsClient,
		config:    config,
	}

	return boggart, nil
}

func (boggart *Boggart) Serve(addr string) error {
	sshHandler := func(s ssh.Session) {
		request, err := ParseRequest(s)
		if err != nil {
			ErrorResponse(err).Write(s)
			_ = s.Exit(1)
			return
		}

		response := boggart.handle(s, s.PublicKey(), request)
		response.Write(s)
		if !response.Success {
			_ = s.Exit(1)
		}
	}

	sshAuthorizer := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return boggart.authorize(ctx, key)
	})

	sshServer := &ssh.Server{
		Addr:    addr,
		Handler: sshHandler,
	}
	if err := sshServer.SetOption(sshAuthorizer); err != nil {
		return fmt.Errorf("failed to set SSH authorizer: %w", err)
	}

	if boggart.config.HostKey != nil {
		hostKey := ssh.HostKeyFile(*boggart.config.HostKey)
		if err := sshServer.SetOption(hostKey); err != nil {
			return fmt.Errorf("failed to set host key: %w", err)
		}
		LogWithoutContext("service.set_host_key", map[string]interface{}{
			"hostKey": *boggart.config.HostKey,
		})
	}

	LogWithoutContext("service.serve", map[string]interface{}{
		"addr": addr,
	})
	return sshServer.ListenAndServe()
}

func (boggart *Boggart) authorize(ctx ssh.Context, key ssh.PublicKey) bool {
	for _, permission := range boggart.config.Permissions {
		if permission.IsKeyAuthorized(key) {
			LogAuthSuccess(ctx, permission)
			return true
		}
	}
	LogAuthFail(ctx)
	return false
}

func (boggart *Boggart) handle(s ssh.Session, key ssh.PublicKey, request *Request) *Response {
	for _, permission := range boggart.config.Permissions {
		if permission.AllowsAssumeRole(key, request.RoleArn) {
			credentials, err := boggart.assumeRole(request)
			if err != nil {
				LogAssumeFail(s, permission, request, err)
				return ErrorResponse(err)
			} else if credentials == nil {
				continue
			} else {
				LogAssumeSuccess(s, permission, request)
				return SuccessResponse(credentials)
			}
		}
	}
	LogAssumeDeny(s, request)
	return ErrorResponse(fmt.Errorf("not allowed to assume role %s", request.RoleArn))
}

func (boggart *Boggart) assumeRole(request *Request) (*sts.Credentials, error) {
	var duration int64 = 900
	input := &sts.AssumeRoleInput{
		RoleArn:         &request.RoleArn,
		RoleSessionName: &request.SessionName,
		DurationSeconds: &duration,
	}
	output, err := boggart.stsClient.AssumeRole(input)
	if err != nil {
		return nil, err
	}
	return output.Credentials, nil
}
