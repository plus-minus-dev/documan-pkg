package contract

// Summary — payload_summary: итоговые суммы документа.
// Числа хранятся как строки с точкой-разделителем ("12345.67").
type Summary struct {
	Amount        string `json:"amount,omitempty"`          // сумма без НДС
	Vat           string `json:"vat,omitempty"`             // сумма НДС
	AmountWithVat string `json:"amount_with_vat,omitempty"` // сумма с НДС (итого к оплате)
}
