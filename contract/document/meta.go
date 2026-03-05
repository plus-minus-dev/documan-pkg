package document

// Meta — payload_meta: служебная информация о документе и его источнике.
type Meta struct {
	// Источник
	Source          string `json:"source,omitempty"`            // "ingest", "connector-ms", "connector-1c"
	SourceAccountID string `json:"source_account_id,omitempty"` // opaque account id источника, не обязательно UUID

	// Номер/дата в системе-источнике (ERP internal name/moment)
	SourceNumber string `json:"source_number,omitempty"`
	SourceDate   string `json:"source_date,omitempty"` // YYYY-MM-DD

	// Состояние в системе-источнике.
	// Приоритет вычисления: deleted -> archived -> active.
	SourceCreatedAt string `json:"source_created_at,omitempty"` // RFC3339 UTC
	SourceUpdatedAt string `json:"source_updated_at,omitempty"` // RFC3339 UTC
	SourceDeletedAt string `json:"source_deleted_at,omitempty"` // RFC3339 UTC
	SourceStatus    string `json:"source_status,omitempty"`     // active | archived | deleted

	// Файл-оригинал (ingest)
	FileName       string `json:"file_name,omitempty"`
	FileExt        string `json:"file_ext,omitempty"`
	FileHash       string `json:"file_hash,omitempty"`        // SHA256
	FileModifiedAt string `json:"file_modified_at,omitempty"` // RFC3339 UTC, если известно

	// S3
	S3Bucket string `json:"s3_bucket,omitempty"`
	S3Key    string `json:"s3_key,omitempty"`

	// Парсинг
	ParseRule  string `json:"parse_rule,omitempty"` // e.g. "invoice_xlsx", "upd_pdf", "ms_supply"
	ParseError string `json:"parse_error,omitempty"`

	// Служебные даты
	CreatedAt   string `json:"created_at,omitempty"`   // RFC3339 UTC, запись создана
	PayloadedAt string `json:"payloaded_at,omitempty"` // RFC3339 UTC, payload сформирован
}
