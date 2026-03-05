package entity

// Store — склад (из ERP).
type Store struct {
	Meta

	// Идентификация
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`

	// ERP
	SourceCode string `json:"source_code,omitempty"`
}
