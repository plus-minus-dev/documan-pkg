package document

// Header — payload_header: реквизиты документа.
// Один struct для всех типов документов (order, shipment, invoice, payment).
// Пустая строка = поле не распознано / не применимо для данного типа.
//
// Универсальная плоская модель:
// - ingest обычно заполняет business fields;
// - ERP connectors могут заполнять только ref fields (*_id) и оставлять business fields пустыми;
// - core enrich-ит business fields по refs при необходимости.
type Header struct {
	// Документ
	Type           string `json:"doc_type,omitempty"`        // оригинальный тип из источника (supply, УПД, счет...)
	Number         string `json:"doc_number,omitempty"`      // номер документа
	Date           string `json:"doc_date,omitempty"`        // дата документа, YYYY-MM-DD
	IncomingNumber string `json:"incoming_number,omitempty"` // входящий номер
	IncomingDate   string `json:"incoming_date,omitempty"`   // входящая дата, YYYY-MM-DD

	// Reference fields (ERP/source refs)
	SellerID string `json:"seller_id,omitempty"` // source entity id продавца
	BuyerID  string `json:"buyer_id,omitempty"`  // source entity id покупателя
	StoreID  string `json:"store_id,omitempty"`  // source entity id склада

	// Продавец (business fields)
	SellerName        string `json:"seller_name,omitempty"`
	SellerINN         string `json:"seller_inn,omitempty"`
	SellerKPP         string `json:"seller_kpp,omitempty"`
	SellerAddress     string `json:"seller_address,omitempty"`
	SellerAccountID   string `json:"seller_account_id,omitempty"`    // расчётный счёт
	SellerBankName    string `json:"seller_bank_name,omitempty"`     // наименование банка
	SellerBankBIK     string `json:"seller_bank_bik,omitempty"`      // БИК банка
	SellerBankCorrAcc string `json:"seller_bank_corr_acc,omitempty"` // корр. счёт банка
	SellerBankAccount string `json:"seller_bank_account,omitempty"`  // расчётный счёт в банке

	// Покупатель (business fields)
	BuyerName        string `json:"buyer_name,omitempty"`
	BuyerINN         string `json:"buyer_inn,omitempty"`
	BuyerKPP         string `json:"buyer_kpp,omitempty"`
	BuyerAddress     string `json:"buyer_address,omitempty"`
	BuyerAccountID   string `json:"buyer_account_id,omitempty"`    // расчётный счёт
	BuyerBankName    string `json:"buyer_bank_name,omitempty"`     // наименование банка
	BuyerBankBIK     string `json:"buyer_bank_bik,omitempty"`      // БИК банка
	BuyerBankCorrAcc string `json:"buyer_bank_corr_acc,omitempty"` // корр. счёт банка
	BuyerBankAccount string `json:"buyer_bank_account,omitempty"`  // расчётный счёт в банке

	// Грузоотправитель (УПД: ФНС ГрузОт)
	// Поля могут быть заполнены частично: например, только Name, если источник не позволяет
	// надежно выделить ИНН/КПП/адрес.
	ShipperName    string `json:"shipper_name,omitempty"`
	ShipperINN     string `json:"shipper_inn,omitempty"`
	ShipperKPP     string `json:"shipper_kpp,omitempty"`
	ShipperAddress string `json:"shipper_address,omitempty"`

	// Грузополучатель (УПД: ФНС ГрузПолуч)
	// Поля могут быть заполнены частично: например, только Name, если источник не позволяет
	// надежно выделить ИНН/КПП/адрес.
	ReceiverName    string `json:"receiver_name,omitempty"`
	ReceiverINN     string `json:"receiver_inn,omitempty"`
	ReceiverKPP     string `json:"receiver_kpp,omitempty"`
	ReceiverAddress string `json:"receiver_address,omitempty"`

	// Склад (business field)
	StoreName    string `json:"store_name,omitempty"`
	StoreAddress string `json:"store_address,omitempty"`

	// Валюта
	CurrencyCode string `json:"currency_code,omitempty"` // ISO 4217 numeric code, e.g. "643"
	CurrencyName string `json:"currency_name,omitempty"`
}
