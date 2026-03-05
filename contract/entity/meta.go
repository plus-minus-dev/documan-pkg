package entity

// Meta — общие метаполя для всех entity (справочников).
// Встраивается в entity structs через embedding, маршалится плоско в JSON.
type Meta struct {
	Source          string `json:"source,omitempty"`             // "connector-ms", "connector-1c"
	SourceAccountID string `json:"source_account_id,omitempty"` // UUID аккаунта ERP
	SourceStatus    string `json:"source_status,omitempty"`     // active | archived | deleted
	SourceCreatedAt string `json:"source_created_at,omitempty"`
	SourceUpdatedAt string `json:"source_updated_at,omitempty"`
	SourceDeletedAt string `json:"source_deleted_at,omitempty"`
	PayloadedAt     string `json:"payloaded_at,omitempty"`
}
