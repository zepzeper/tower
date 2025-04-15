package data

// OrderTax represents an order-level tax
type OrderTax struct {
	Code        string   `json:"code"`
	Rate        *float64 `json:"rate"`
	TotalAmount float64  `json:"total_amount"`
}

// TaxWithBreakdown represents a tax with amount breakdown
type TaxWithBreakdown struct {
	Amount             float64          `json:"amount"`
	AmountBreakdown    *AmountBreakdown `json:"amount_breakdown,omitempty"`
	Code               string           `json:"code"`
	PurchaseTax        PurchaseTax      `json:"purchase_tax"`
	Rate               float64          `json:"rate,omitempty"`
	TaxCalculationRule string           `json:"tax_calculation_rule,omitempty"`
}

// Tax represents a tax
type Tax struct {
	Amount             float64     `json:"amount"`
	Code               string      `json:"code"`
	OriginUnitAmount   float64     `json:"origin_unit_amount,omitempty"`
	PurchaseTax        PurchaseTax `json:"purchase_tax"`
	Rate               float64     `json:"rate,omitempty"`
	TaxCalculationRule string      `json:"tax_calculation_rule,omitempty"`
}

// PurchaseTax represents purchase tax information
type PurchaseTax struct {
	PurchaseAmount float64 `json:"purchase_amount"`
	PurchaseRate   float64 `json:"purchase_rate,omitempty"`
}

// ShippingTax represents tax on shipping
type ShippingTax struct {
	Amount             float64     `json:"amount"`
	Code               string      `json:"code"`
	PurchaseTax        PurchaseTax `json:"purchase_tax"`
	Rate               float64     `json:"rate,omitempty"`
	TaxCalculationRule string      `json:"tax_calculation_rule,omitempty"`
}

// CommissionTax represents tax on commission
type CommissionTax struct {
	Amount float64 `json:"amount"`
	Code   string  `json:"code"`
	Rate   float64 `json:"rate"`
}

// CommissionTaxShort represents simplified commission tax information
type CommissionTaxShort struct {
	Amount float64 `json:"amount"`
	Code   string  `json:"code"`
}
