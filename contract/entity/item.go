package entity

// Item — товар, услуга, комплект или модификация/вариант из ERP.
type Item struct {
	Meta

	// Идентификация
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"item_type,omitempty"` // "product" | "service" | "bundle" | "variant"
	Variant  string   `json:"variant,omitempty"`   // характеристика / модификация / вариант
	Article  string   `json:"article,omitempty"`   // артикул / article
	ItemCode string   `json:"item_code,omitempty"` // код товара из ERP/source business code; не article
	GTIN     []string `json:"gtin,omitempty"`      // ГТИН (массив, с карточки)
	Barcodes []string `json:"barcodes,omitempty"`  // прочие штрихкоды
}
