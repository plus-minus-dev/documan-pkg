package entity

// Party — контрагент или организация.
// Собирательная сущность: продавцы, покупатели, ООО, ИП.
// Роль (продавец/покупатель) определяется позицией в документе, не свойством Party.
type Party struct {
	Meta

	// Идентификация
	Name        string `json:"name,omitempty"`        // краткое название
	LegalTitle  string `json:"legal_title,omitempty"` // полное юридическое название
	INN         string `json:"inn,omitempty"`
	KPP         string `json:"kpp,omitempty"`
	CompanyType string `json:"companyType,omitempty"`
	Address     string `json:"address,omitempty"`

	// ERP/source business code
	SourceCode string `json:"source_code,omitempty"`
}
