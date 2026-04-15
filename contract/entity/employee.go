package entity

type Employee struct {
	Meta

	LastName   string `json:"last_name,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
	FullName   string `json:"full_name,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Position   string `json:"position,omitempty"`
	SourceCode string `json:"source_code,omitempty"`
}
