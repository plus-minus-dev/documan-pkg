package contract

// Position — одна строка (позиция) товара/услуги в документе.
// Объединяет поля УПД (16 колонок) и Счетов (10 колонок).
// Суммы хранятся в копейках (int). Количество — float64.
type Position struct {
	LineNo   string  `json:"line_no,omitempty"`   // порядковый номер строки
	Article  string  `json:"article,omitempty"`   // артикул / код товара
	Name     string  `json:"name,omitempty"`      // наименование товара / услуги
	Unit     string  `json:"unit,omitempty"`      // единица измерения (текст)
	UnitCode string  `json:"unit_code,omitempty"` // код единицы измерения (ОКЕИ)
	Quantity float64 `json:"quantity,omitempty"`   // количество

	Price        int `json:"price,omitempty"`          // цена за единицу без НДС (копейки)
	PriceWithVat int `json:"price_with_vat,omitempty"` // цена за единицу с НДС (копейки)
	Amount       int `json:"amount,omitempty"`          // сумма без НДС (копейки)
	Excise       int `json:"excise,omitempty"`          // акциз (копейки)

	VatRate       string `json:"vat_rate,omitempty"`        // ставка НДС (%, "без НДС")
	Vat           int    `json:"vat,omitempty"`             // сумма НДС (копейки)
	AmountWithVat int    `json:"amount_with_vat,omitempty"` // сумма с НДС (копейки)

	// Таможенные поля (УПД)
	CommodityCode string `json:"commodity_code,omitempty"` // код вида товара
	CountryCode   string `json:"country_code,omitempty"`   // цифровой код страны
	CountryName   string `json:"country_name,omitempty"`   // краткое название страны
	CustomsGTD    string `json:"customs_gtd,omitempty"`    // номер ГТД
}
