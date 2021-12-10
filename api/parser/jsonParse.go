package parser

import (
	"encoding/json"
	"io"
)

const (
	READ_PAYLOAD_CHUNK_LENGTH int64 = 256
)

// utility to parse the bytes of json format to Struct
func Json(body io.ReadCloser, req interface{}) error {
	defer body.Close()
	err := json.NewDecoder(body).Decode(req)
	if err != nil {
		return err
	}
	return nil
}
