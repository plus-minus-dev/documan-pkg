package entity

type Project struct {
	Meta

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SourceCode  string `json:"source_code,omitempty"`
}
