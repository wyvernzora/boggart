package boggart

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/wyvernzora/boggart/pkg/config"
	"github.com/wyvernzora/boggart/pkg/model"
	"golang.org/x/crypto/ssh"
)

const DefaultAssumeRoleDuration int64 = 900

type Boggart struct {
	stsClient   *sts.STS
	permissions []config.Permission
}

func New(config *config.Config) *Boggart {
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region:      &config.AwsConfig.Region,
		Credentials: config.AwsConfig.Credentials,
	}))
	stsClient := sts.New(awsSession)

	return &Boggart{
		stsClient:   stsClient,
		permissions: config.Permissions,
	}
}

func (b *Boggart) Authenticate(ctx *Context, key ssh.PublicKey) bool {
	for _, p := range b.permissions {
		if p.IsKeyAuthorized(key) {
			ctx.Authenticated(&p)
			ctx.Logger.Info().Msg("authentication success")
			return true
		}
	}
	ctx.Logger.Info().Msg("authentication fail")
	return false
}

func (b *Boggart) Authorize(ctx *Context, key ssh.PublicKey, req *model.Request) bool {
	for _, p := range b.permissions {
		if p.IsKeyAuthorized(key) && p.CanAssumeRole(req.RoleArn) {
			ctx.Authorized(&p)
			ctx.Logger.Info().Msg("authorization success")
			return true
		}
	}
	ctx.Logger.Info().Msg("authorization fail")
	return false
}

func (b *Boggart) AssumeRole(ctx *Context, req *model.Request) *model.Response {
	var duration = DefaultAssumeRoleDuration

	if req.SessionName == "" {
		req.SessionName = fmt.Sprintf("%s@%s", ctx.Permission.Name, ctx.RemoteAddr)
	}
	ctx.Request(req)

	input := &sts.AssumeRoleInput{
		RoleArn:         &req.RoleArn,
		RoleSessionName: &req.SessionName,
		DurationSeconds: &duration,
	}
	output, err := b.stsClient.AssumeRole(input)
	if err != nil {
		ctx.Logger.Error().
			Err(err).
			Msg("assume role fail")
		return model.NewErrorResponse(err)
	}
	ctx.Logger.Info().Msg("assume role success")
	return model.NewSuccessResponse(output.Credentials)
}
