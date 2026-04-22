package document

// Summary — payload_summary: итоговые суммы и агрегаты документа.
// Денежные поля — *int64 (копейки): nil = нет данных, 0 = документ без сумм.
// Qty/SKU/Lines намеренно сериализуются всегда, даже если документ не содержит позиций.
type Summary struct {
	// Суммы (копейки)
	Amount        *int64 `json:"amount,omitempty"`          // сумма без НДС
	TaxAmount     *int64 `json:"tax_amount,omitempty"`      // сумма НДС
	AmountWithTax *int64 `json:"amount_with_tax,omitempty"` // сумма с НДС
	Discount      *int64 `json:"discount,omitempty"`        // скидка, копейки

	// Агрегаты
	// Qty: decimal string, dot separator, no thousand separators, normalized.
	// Examples: "1", "1.25", "0.500".
	Qty   string `json:"qty"`
	SKU   int    `json:"sku"`   // количество уникальных артикулов
	Lines int    `json:"lines"` // количество строк в документе
}
