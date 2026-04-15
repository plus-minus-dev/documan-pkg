package entity

type Country struct {
	Meta

	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"` // ОКСМ
}
