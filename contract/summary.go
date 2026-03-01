package contract

// Summary — payload_summary: итоговые суммы документа.
// Суммы хранятся в копейках (int). Пример: 1234567 = 12345.67 руб.
type Summary struct {
	Amount        int `json:"amount,omitempty"`          // сумма без НДС (копейки)
	Vat           int `json:"vat,omitempty"`             // сумма НДС (копейки)
	AmountWithVat int `json:"amount_with_vat,omitempty"` // сумма с НДС (копейки)
}
