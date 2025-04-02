package orchestrator

import (
	"encoding/json"
	"io"
)

type StdinReader struct {
	Reader io.Reader
}

func (r *StdinReader) ReadInput() (map[string]interface{}, error) {
	var data map[string]interface{}
	decoder := json.NewDecoder(r.Reader)
	err := decoder.Decode(&data)
	return data, err
}
