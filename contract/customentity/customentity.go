package customentity

import "github.com/plus-minus-dev/documan-pkg/contract/document"

type Element struct {
	Name         string               `json:"name,omitempty"`
	Code         string               `json:"code,omitempty"`
	Description  string               `json:"description,omitempty"`
	ExternalCode string               `json:"external_code,omitempty"`
	MetadataID   string               `json:"metadata_id,omitempty"`
	MetadataName string               `json:"metadata_name,omitempty"`
	Attributes   []document.Attribute `json:"attributes,omitempty"`
}
