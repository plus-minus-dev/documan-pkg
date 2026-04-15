package customentity

type Element struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	MetadataID   string `json:"metadata_id,omitempty"`
	MetadataName string `json:"metadata_name,omitempty"`
	SourceCode   string `json:"source_code,omitempty"`
}
