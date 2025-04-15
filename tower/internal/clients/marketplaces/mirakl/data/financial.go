package data

// Cancelation represents a canceled item
type Cancelation struct {
	Amount                  float64             `json:"amount"`
	AmountBreakdown         AmountBreakdown     `json:"amount_breakdown"`
	CommissionAmount        float64             `json:"commission_amount"`
	CommissionTaxes         []interface{}       `json:"commission_taxes"`
	CommissionTotalAmount   float64             `json:"commission_total_amount"`
	CreatedDate             string              `json:"created_date"`
	EcoContributions        []interface{}       `json:"eco_contributions"`
	Fees                    []Fee               `json:"fees"`
	ID                      string              `json:"id"`
	PurchaseInformation     PurchaseInformation `json:"purchase_information"`
	Quantity                int                 `json:"quantity"`
	ReasonCode              string              `json:"reason_code"`
	ShippingAmount          float64             `json:"shipping_amount"`
	ShippingAmountBreakdown AmountBreakdown     `json:"shipping_amount_breakdown"`
	ShippingTaxes           []TaxWithBreakdown  `json:"shipping_taxes"`
	Taxes                   []TaxWithBreakdown  `json:"taxes"`
}

// Refund represents a refund
type Refund struct {
	Amount                float64              `json:"amount"`
	CommissionAmount      float64              `json:"commission_amount"`
	CommissionTaxAmount   float64              `json:"commission_tax_amount"`
	CommissionTaxes       []CommissionTaxShort `json:"commission_taxes"`
	CommissionTotalAmount float64              `json:"commission_total_amount"`
	CreatedDate           string               `json:"created_date"`
	EcoContributions      []interface{}        `json:"eco_contributions"`
	Fees                  []Fee                `json:"fees"`
	ID                    string               `json:"id"`
	OrderRefundID         string               `json:"order_refund_id"`
	PurchaseInformation   PurchaseInformation  `json:"purchase_information"`
	Quantity              int                  `json:"quantity"`
	ReasonCode            string               `json:"reason_code"`
	ShippingAmount        float64              `json:"shipping_amount"`
	ShippingTaxes         []ShippingTax        `json:"shipping_taxes"`
	State                 string               `json:"state"`
	TaxLegalNotice        string               `json:"tax_legal_notice"`
	Taxes                 []Tax                `json:"taxes"`
}

// AmountBreakdown represents a breakdown of a monetary amount
type AmountBreakdown struct {
	Parts []AmountPart `json:"parts"`
}

// AmountPart represents a part of an amount breakdown
type AmountPart struct {
	Amount                float64 `json:"amount"`
	Commissionable        bool    `json:"commissionable"`
	DebitableFromCustomer bool    `json:"debitable_from_customer"`
	PayableToShop         bool    `json:"payable_to_shop"`
}

// PurchaseInformation represents purchase details
type PurchaseInformation struct {
	PurchaseCommissionOnFees     TotalAmountContainer `json:"purchase_commission_on_fees"`
	PurchaseCommissionOnPrice    float64              `json:"purchase_commission_on_price"`
	PurchaseCommissionOnShipping float64              `json:"purchase_commission_on_shipping"`
	PurchaseFeeAmount            TotalAmountContainer `json:"purchase_fee_amount"`
	PurchasePrice                float64              `json:"purchase_price"`
	PurchaseShippingPrice        float64              `json:"purchase_shipping_price"`
	PurchaseUnitPrice            float64              `json:"purchase_unit_price,omitempty"`
}

// TotalAmountContainer represents a container for total amount
type TotalAmountContainer struct {
	TotalAmount float64 `json:"total_amount"`
}

// Fee represents a fee
type Fee struct {
	Amount float64 `json:"amount"`
	Code   string  `json:"code"`
}

// ListRefunds retrieves all refunds for a given order line
func (o Orders) ListRefunds(orderID string, orderLineID string) ([]Refund, error) {
	order, err := o.Get(orderID)
	if err != nil {
		return nil, err
	}
	
	for _, line := range order.OrderLines {
		if line.OrderLineID == orderLineID {
			return line.Refunds, nil
		}
	}
	
	return nil, nil
}
