package format

import (
	"bufio"
	"bytes"
	"fmt"
	boggart "github.com/wyvernzora/boggart/pkg/model"
)

type ShellFormatter struct {
}

func (f ShellFormatter) Name() string {
	return "shell"
}

func (f ShellFormatter) Apply(response *boggart.Response) ([]byte, error) {
	var data bytes.Buffer
	writer := bufio.NewWriter(&data)

	if !response.Success {
		fmt.Fprintf(writer, "export BOGGART_SUCCESS=0\n")
		fmt.Fprintf(writer, "export BOGGART_ERROR='%s'\n", *response.Error)
	} else {
		fmt.Fprintf(writer, "export BOGGART_SUCCESS=1\n")
		fmt.Fprintf(writer, "export BOGGART_ERROR=\n")

		creds := *response.Credentials
		fmt.Fprintf(writer, "export BOGGART_EXPIRES_AT='%s'\n", creds.ExpiresAt)
		fmt.Fprintf(writer, "export AWS_ACCESS_KEY_ID='%s'\n", creds.AccessKeyId)
		fmt.Fprintf(writer, "export AWS_SECRET_ACCESS_KEY='%s'\n", creds.SecretAccessKey)
		fmt.Fprintf(writer, "export AWS_SESSION_TOKEN='%s'\n", creds.SessionToken)
	}
	writer.Flush()

	return data.Bytes(), nil
}
