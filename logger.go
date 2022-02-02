package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

type Context interface {
	User() string
	RemoteAddr() net.Addr
}

func LogWithoutContext(ev string, values map[string]interface{}) {
	entry := make(map[string]interface{})
	entry["ev"] = ev

	for k, v := range values {
		entry[k] = v
	}
	data, _ := json.Marshal(entry)
	fmt.Println(string(data))
}

func LogWithContext(ctx Context, ev string, values map[string]interface{}) {
	entry := make(map[string]interface{})

	entry["t"] = time.Now().Format(time.RFC3339)
	entry["user"] = ctx.User()
	entry["ip"] = ctx.RemoteAddr().String()

	for k, v := range values {
		entry[k] = v
	}
	LogWithoutContext(ev, entry)
}

func LogAuthSuccess(ctx Context, permission *Permission) {
	LogWithContext(ctx, "authn.success", map[string]interface{}{
		"key": permission.Name,
	})
}

func LogAuthFail(ctx Context) {
	LogWithContext(ctx, "authn.fail", map[string]interface{}{})
}

func LogAssumeSuccess(ctx Context, permission *Permission, request *Request) {
	LogWithContext(ctx, "assume.success", map[string]interface{}{
		"key":         permission.Name,
		"roleArn":     request.RoleArn,
		"sessionName": request.SessionName,
	})
}

func LogAssumeFail(ctx Context, permission *Permission, request *Request, err error) {
	LogWithContext(ctx, "assume.fail", map[string]interface{}{
		"key":         permission.Name,
		"roleArn":     request.RoleArn,
		"sessionName": request.SessionName,
		"error":       strings.ReplaceAll(err.Error(), "\n", " "),
	})
}

func LogAssumeDeny(ctx Context, request *Request) {
	LogWithContext(ctx, "assume.deny", map[string]interface{}{
		"roleArn":     request.RoleArn,
		"sessionName": request.SessionName,
	})
}
