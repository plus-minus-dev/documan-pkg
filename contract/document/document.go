package document

// Canonical document types.
const (
	TypeOrder    = "order"
	TypeShipment = "shipment"
	TypeInvoice  = "invoice"
	TypePayment  = "payment"
)

// Document direction.
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
