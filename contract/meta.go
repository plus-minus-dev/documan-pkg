package contract

// Meta — payload_meta: общие сведения о документе.
type Meta struct {
	DocType string `json:"doc_type"` // "УПД", "счет", "счет-оферта", "счет-спецификация", "UNKNOWN"
	Format  string `json:"format"`   // "xls", "xlsx", "pdf"
}

// Known DocType values.
const (
	DocTypeUPD              = "УПД"
	DocTypeInvoice          = "счет"
	DocTypeInvoiceOffer     = "счет-оферта"
	DocTypeInvoiceSpec      = "счет-спецификация"
	DocTypeUnknown          = "UNKNOWN"
)

// Known Format values.
const (
	FormatXLS  = "xls"
	FormatXLSX = "xlsx"
	FormatPDF  = "pdf"
)
