package models

// Beer
type Beer struct {
	Model
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	FileID      *uint  `json:"file_id,omitempty"`
	File        *File  `json:"file,omitempty"`
}
