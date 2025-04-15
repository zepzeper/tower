package data

// InvoiceDetails represents invoice information
type InvoiceDetails struct {
	DocumentDetails []DocumentDetail `json:"document_details"`
	PaymentTerms    PaymentTerms     `json:"payment_terms"`
}

// DocumentDetail represents document format details
type DocumentDetail struct {
	Format string `json:"format"`
}

// PaymentTerms represents payment term details
type PaymentTerms struct {
	Days int    `json:"days"`
	Type string `json:"type"`
}
