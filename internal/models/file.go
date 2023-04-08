package models

import (
	"beer-api/internal/core/config"
	"encoding/json"
	"fmt"
)

// File file model
type File struct {
	Model
	Name string `json:"name"`
}

// MarshalJSON custom image json
func (f File) MarshalJSON() ([]byte, error) {
	type Alias File
	fileModel := &struct {
		*Alias
		FileURL string `json:"file_url,omitempty"`
	}{
		Alias:   (*Alias)(&f),
		FileURL: fmt.Sprintf("%s/api/v1/public/%s", config.CF.App.URL, f.Name),
	}

	return json.Marshal(fileModel)
}
