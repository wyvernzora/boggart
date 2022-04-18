package format

import (
	"encoding/json"
	boggart "github.com/wyvernzora/boggart/pkg/model"
)

type JsonFormatter struct {
}

func (f JsonFormatter) Name() string {
	return "json"
}

func (f JsonFormatter) Apply(response *boggart.Response) ([]byte, error) {
	return json.Marshal(response)
}
