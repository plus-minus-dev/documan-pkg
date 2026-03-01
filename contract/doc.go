// Package contract defines canonical field names and types for document payloads
// exchanged between DocuMan services (Ingest → Core → Frontend).
//
// All 4 payload types (Meta, Header, Summary, Positions) are transmitted as JSON []byte
// over gRPC. This package provides typed structs and constants so every service
// serialises/deserialises with the same field names.
//
// Usage:
//
//	// Producer (Ingest):
//	h := contract.Header{SellerName: "ООО Ромашка", SellerINN: "7733456789"}
//	data, _ := json.Marshal(h)
//
//	// Consumer (Core / Frontend):
//	var h contract.Header
//	json.Unmarshal(payloadHeader, &h)
package contract
