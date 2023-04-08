package utils

import (
	"bytes"
	"encoding/json"
	"os"
)

// PrintFormatJSON print format json
func PrintFormatJSON(n interface{}) {
	b, _ := json.MarshalIndent(n, "", "\t")
	_, _ = os.Stdout.Write(b)
}

// ReadJSONFile read json file
func ReadJSONFile(path string, entities interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	body := bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	err = json.Unmarshal(body, entities)
	if err != nil {
		return err
	}

	return nil
}
