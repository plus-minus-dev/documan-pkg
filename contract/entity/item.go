package entity

// Item — товар, услуга, комплект (карточка из ERP).
// Название предварительное — будет пересмотрено при реализации товарного matching.
type Item struct {
	Meta

	// Идентификация
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"item_type,omitempty"` // "product" | "service" | "bundle" | "variant"
	SKU      string   `json:"sku,omitempty"`
	GTIN     []string `json:"gtin,omitempty"`     // ГТИН (массив, с карточки)
	Barcodes []string `json:"barcodes,omitempty"` // прочие штрихкоды

	// ERP
	SourceCode string `json:"source_code,omitempty"`
}
