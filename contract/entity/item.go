package entity

import "github.com/plus-minus-dev/documan-pkg/contract/document"

type Item struct {
	Meta

	Name            string               `json:"name,omitempty"`
	Type            string               `json:"item_type,omitempty"`
	Description     string               `json:"description,omitempty"`
	Article         string               `json:"article,omitempty"`
	ItemCode        string               `json:"item_code,omitempty"`
	GTIN            []string             `json:"gtin,omitempty"`
	Barcodes        []string             `json:"barcodes,omitempty"`
	Characteristics []string             `json:"characteristics,omitempty"`
	ProductID       string               `json:"product_id,omitempty"`
	Attributes      []document.Attribute `json:"attributes,omitempty"`
}
