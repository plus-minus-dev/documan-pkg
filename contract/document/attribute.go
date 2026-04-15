package document

type Attribute struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	StringVal string   `json:"string_val,omitempty"`
	LongVal   *int64   `json:"long_val,omitempty"`
	DoubleVal *float64 `json:"double_val,omitempty"`
	BoolVal   *bool    `json:"bool_val,omitempty"`
	RefID     string   `json:"ref_id,omitempty"`
	RefName   string   `json:"ref_name,omitempty"`
}
