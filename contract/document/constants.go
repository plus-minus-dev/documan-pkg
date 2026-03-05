package document

// Canonical document types for the classification layer in core.
// These constants are not values for payload_header.doc_type.
const (
	TypeOrder    = "order"
	TypeShipment = "shipment"
	TypeInvoice  = "invoice"
	TypePayment  = "payment"
)

// Canonical document direction for the classification layer in core.
const (
	DirectionForward = "forward"
	DirectionReverse = "reverse"
)

// Source status in origin system.
const (
	SourceStatusActive   = "active"
	SourceStatusArchived = "archived"
	SourceStatusDeleted  = "deleted"
)
