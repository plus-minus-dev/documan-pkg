package document

// Position — одна строка (позиция) товара/услуги в документе.
// Денежные поля — *int64 (копейки): nil = нет данных, 0 = валидный ноль.
// Количества и процентные значения — decimal string (dot separator, no thousand separators).
//
// Универсальная плоская модель:
// - ingest обычно заполняет business fields;
// - ERP connectors могут заполнять только item_id и оставлять business fields пустыми;
// - core enrich-ит business fields по item_id при необходимости.
type Position struct {
	// Идентификация строки
	LineNumber int    `json:"line_number"`          // нормализованный 1-based индекс строки в payload
	LineLabel  string `json:"line_label,omitempty"` // исходное значение номера строки из документа: "1", "1а", "I" ...
	ItemType   string `json:"item_type,omitempty"`  // "product" | "service" | "work" | "rights"

	// Reference field (ERP/source ref)
	ItemID string `json:"item_id,omitempty"` // source entity id товара/услуги в ERP

	// Товар/услуга (business fields)
	Name     string   `json:"name,omitempty"`      // наименование
	Variant  string   `json:"variant,omitempty"`   // характеристика
	Article  string   `json:"article,omitempty"`   // артикул / article
	ItemCode string   `json:"item_code,omitempty"` // код товара/работ/услуг из документа или источника; не article, не gtin, не hs_code
	GTIN     string   `json:"gtin,omitempty"`      // ГТИН из документа (ФНС)
	Barcodes []string `json:"barcodes,omitempty"`  // прочие штрихкоды (code128, UPC, EAN-8...)

	// Количество и единицы
	// Quantity: decimal string, dot separator, no thousand separators, normalized.
	// Examples: "1", "1.25", "0.500".
	Quantity string `json:"quantity,omitempty"`
	UOMCode  string `json:"uom_code,omitempty"` // код ОКЕИ
	UOMName  string `json:"uom_name,omitempty"` // ед. изм. текст

	// Цена за единицу (копейки)
	UnitPrice        *int64 `json:"unit_price,omitempty"`          // без НДС
	UnitPriceWithTax *int64 `json:"unit_price_with_tax,omitempty"` // с НДС

	// Суммы по строке (копейки)
	Amount        *int64 `json:"amount,omitempty"`          // без НДС
	AmountWithTax *int64 `json:"amount_with_tax,omitempty"` // с НДС

	// Налог
	TaxRate   string `json:"tax_rate,omitempty"`   // normalized string, e.g. "20%", "10%", "без НДС"
	TaxAmount *int64 `json:"tax_amount,omitempty"` // сумма НДС

	// Скидка
	// DiscountRate: decimal string, e.g. "5", "10.5".
	// Discount заполняется только если producer может извлечь или вычислить сумму скидки без неоднозначности.
	DiscountRate string `json:"discount_rate,omitempty"`
	Discount     *int64 `json:"discount,omitempty"` // скидка, копейки

	// Акциз (копейки, 0 = "без акциза")
	Excise *int64 `json:"excise,omitempty"`

	// Страна и таможня
	HSCode      string `json:"hs_code,omitempty"`      // код ТН ВЭД
	CountryCode string `json:"country_code,omitempty"` // код страны ОКСМ
	CountryName string `json:"country_name,omitempty"` // название страны
	CustomsGTD  string `json:"customs_gtd,omitempty"`  // номер ГТД

	// Маркировка
	TrackingType        string   `json:"tracking_type,omitempty"`        // MILK, TOBACCO, ...
	IdentificationCodes []string `json:"identification_codes,omitempty"` // коды маркировки
}
