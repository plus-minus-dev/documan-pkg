package entity

// Meta — общие метаполя для всех entity (справочников).
// Встраивается в entity structs через embedding, маршалится плоско в JSON.
type Meta struct {
	Source          string `json:"source,omitempty"`            // "connector-ms", "connector-1c"
	SourceAccountID string `json:"source_account_id,omitempty"` // opaque account id источника, не обязательно UUID
	SourceStatus    string `json:"source_status,omitempty"`     // active | archived | deleted; приоритет: deleted -> archived -> active
	SourceCreatedAt string `json:"source_created_at,omitempty"` // RFC3339 UTC
	SourceUpdatedAt string `json:"source_updated_at,omitempty"` // RFC3339 UTC
	SourceDeletedAt string `json:"source_deleted_at,omitempty"` // RFC3339 UTC
	PayloadedAt     string `json:"payloaded_at,omitempty"`      // RFC3339 UTC
}
