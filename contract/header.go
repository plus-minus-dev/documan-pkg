package contract

// Header — payload_header: реквизиты документа.
// Все строковые поля; пустая строка = не распознано.
type Header struct {
	// Документ
	DocType   string `json:"doc_type,omitempty"`   // дублирует meta для удобства
	DocNumber string `json:"doc_number,omitempty"` // номер документа
	DocDate   string `json:"doc_date,omitempty"`   // дата документа (DD.MM.YYYY)

	// Продавец / Поставщик
	SellerName    string `json:"seller_name,omitempty"`
	SellerAddress string `json:"seller_address,omitempty"`
	SellerINN     string `json:"seller_inn,omitempty"`
	SellerKPP     string `json:"seller_kpp,omitempty"`

	// Покупатель
	BuyerName    string `json:"buyer_name,omitempty"`
	BuyerAddress string `json:"buyer_address,omitempty"`
	BuyerINN     string `json:"buyer_inn,omitempty"`
	BuyerKPP     string `json:"buyer_kpp,omitempty"`

	// Грузоотправитель / Грузополучатель (УПД)
	ShipperName  string `json:"shipper_name,omitempty"`
	ReceiverName string `json:"receiver_name,omitempty"`

	// Валюта
	CurrencyName string `json:"currency_name,omitempty"`
	CurrencyCode string `json:"currency_code,omitempty"` // ISO 4217, e.g. "643"

	// Ошибка парсинга (если формат не распознан)
	ParseError string `json:"parse_error,omitempty"`
}
