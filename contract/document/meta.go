package document

// Meta — payload_meta: служебная информация о документе и его источнике.
type Meta struct {
	// Источник
	Source string `json:"source,omitempty"` // для original: "ms_files", "ms_dnd", "web" или "my_dnd" или "ui_dnd". Для erp: "connector-ms", "connector-1c"

	SourceAccountID string     `json:"source_account_id,omitempty"` // opaque account id источника, не обязательно UUID
	Files           []MetaFile `json:"files,omitempty"`              // привязка к ERP-документам (source=ms_files)

	// Номер/дата в системе-источнике (ERP internal name/moment)
	SourceNumber string `json:"source_number,omitempty"`
	SourceDate   string `json:"source_date,omitempty"` // YYYY-MM-DD

	// Состояние в системе-источнике
	SourceCreatedAt string `json:"source_created_at,omitempty"` // RFC3339 UTC
	SourceUpdatedAt string `json:"source_updated_at,omitempty"` // RFC3339 UTC
	SourceDeletedAt string `json:"source_deleted_at,omitempty"` // RFC3339 UTC
	SourceState       string `json:"source_state,omitempty"`       // статус документа в ERP
	SourceRate        string `json:"source_rate,omitempty"`        // курс валюты в источнике
	SourceDescription string `json:"source_description,omitempty"` // описание / примечание к документу из источника

	// Файл-оригинал (ingest)
	FileName       string `json:"file_name,omitempty"`
	FileExt        string `json:"file_ext,omitempty"`
	FileHash       string `json:"file_hash,omitempty"`        // SHA256
	FileSize       int64  `json:"file_size,omitempty"`        // размер файла в байтах
	FilePages      int    `json:"file_pages,omitempty"`       // количество страниц (PDF) или листов (XLSX/XLS)
	FileMime       string `json:"file_mime,omitempty"`        // MIME type (application/pdf, application/vnd.ms-excel, ...)
	FileModifiedAt string `json:"file_modified_at,omitempty"` // RFC3339 UTC, если известно

	// S3
	S3Bucket string `json:"s3_bucket,omitempty"`
	S3Key    string `json:"s3_key,omitempty"`

	// Парсинг
	ParseRule  string `json:"parse_rule,omitempty"` // e.g. "invoice_xlsx", "upd_pdf", "ms_supply"
	ParseError string `json:"parse_error,omitempty"`

	// Служебные даты и состояние записи в Documan
	CreatedAt   string `json:"created_at,omitempty"`   // RFC3339 UTC, запись создана
	UpdatedAt   string `json:"updated_at,omitempty"`   // RFC3339 UTC, запись обновлена
	DeletedAt   string `json:"deleted_at,omitempty"`   // RFC3339 UTC, запись удалена (soft delete)
	PayloadedAt string `json:"payloaded_at,omitempty"` // RFC3339 UTC, payload сформирован
	Archived    bool   `json:"archived,omitempty"`     // soft delete флаг записи в Documan

	// Трассировка (закомментировано — раскомментировать при необходимости)
	// ConnectionType         string `json:"connection_type,omitempty"`          // auth type соединения (json_api, vendor_app, ...)
	// LastSourceConnectionID string `json:"last_source_connection_id,omitempty"` // UUID connection, который последним обновил запись
}

// MetaFile — привязка original-документа к ERP-документу по entity_type/entity_id.
type MetaFile struct {
	EntityType string `json:"entity_type,omitempty"`
	EntityID   string `json:"entity_id,omitempty"`
}
