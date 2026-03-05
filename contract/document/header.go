package document

// Header — payload_header: реквизиты документа.
// Один struct для всех типов документов (order, shipment, invoice, payment).
// Пустая строка = поле не распознано / не применимо для данного типа.
type Header struct {
	// Документ
	Type   string `json:"doc_type,omitempty"`   // оригинальный тип из источника (supply, УПД, счет...)
	Number string `json:"doc_number,omitempty"` // номер документа
	Date   string `json:"doc_date,omitempty"`   // дата документа

	// Продавец
	SellerName    string `json:"seller_name,omitempty"`
	SellerINN     string `json:"seller_inn,omitempty"`
	SellerKPP     string `json:"seller_kpp,omitempty"`
	SellerAddress string `json:"seller_address,omitempty"`

	// Покупатель
	BuyerName    string `json:"buyer_name,omitempty"`
	BuyerINN     string `json:"buyer_inn,omitempty"`
	BuyerKPP     string `json:"buyer_kpp,omitempty"`
	BuyerAddress string `json:"buyer_address,omitempty"`

	// Грузоотправитель (УПД: ФНС ГрузОт)
	ShipperName    string `json:"shipper_name,omitempty"`
	ShipperINN     string `json:"shipper_inn,omitempty"`
	ShipperKPP     string `json:"shipper_kpp,omitempty"`
	ShipperAddress string `json:"shipper_address,omitempty"`

	// Грузополучатель (УПД: ФНС ГрузПолуч)
	ReceiverName    string `json:"receiver_name,omitempty"`
	ReceiverINN     string `json:"receiver_inn,omitempty"`
	ReceiverKPP     string `json:"receiver_kpp,omitempty"`
	ReceiverAddress string `json:"receiver_address,omitempty"`

	// Склад
	StoreName string `json:"store_name,omitempty"`

	// Валюта
	CurrencyCode string `json:"currency_code,omitempty"` // ISO 4217, e.g. "643"
	CurrencyName string `json:"currency_name,omitempty"`
}
