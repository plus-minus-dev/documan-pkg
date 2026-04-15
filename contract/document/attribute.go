package document

type Attribute struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Value   string `json:"value"`
	ValueID string `json:"value_id,omitempty"`
	DictID  string `json:"dict_id,omitempty"`
}
