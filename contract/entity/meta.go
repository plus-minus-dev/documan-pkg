package entity

// Meta — общие метаполя для всех entity (справочников).
// Встраивается в entity structs через embedding, маршалится плоско в JSON.
type Meta struct {
	// Источник
	Source          string `json:"source,omitempty"`            // "connector-ms", "connector-1c"
	SourceAccountID string `json:"source_account_id,omitempty"` // opaque account id источника, не обязательно UUID

	// Состояние в системе-источнике
	SourceCreatedAt string `json:"source_created_at,omitempty"` // RFC3339 UTC
	SourceUpdatedAt string `json:"source_updated_at,omitempty"` // RFC3339 UTC
	SourceDeletedAt string `json:"source_deleted_at,omitempty"` // RFC3339 UTC
	SourceArchived  bool   `json:"source_archived,omitempty"`   // заархивирован в ERP

	// Служебные даты и состояние записи в Documan
	CreatedAt   string `json:"created_at,omitempty"`   // RFC3339 UTC, запись создана
	UpdatedAt   string `json:"updated_at,omitempty"`   // RFC3339 UTC, запись обновлена
	DeletedAt   string `json:"deleted_at,omitempty"`   // RFC3339 UTC, запись удалена (soft delete)
	PayloadedAt string `json:"payloaded_at,omitempty"` // RFC3339 UTC, payload сформирован
	Archived    bool   `json:"archived,omitempty"`     // soft delete флаг записи в Documan
}
