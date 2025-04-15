package data

import (
	"net/url"
)

// Order represents a single order with all its details
type Order struct {
	AcceptanceDecisionDate      string              `json:"acceptance_decision_date"`
	CanCancel                   bool                `json:"can_cancel"`
	CanShopShip                 bool                `json:"can_shop_ship"`
	Channel                     Channel             `json:"channel"`
	CommercialID                string              `json:"commercial_id"`
	CreatedDate                 string              `json:"created_date"`
	CurrencyIsoCode             string              `json:"currency_iso_code"`
	Customer                    Customer            `json:"customer"`
	CustomerDebitedDate         string              `json:"customer_debited_date"`
	CustomerDirectlyPaysSeller  bool                `json:"customer_directly_pays_seller"`
	CustomerNotificationEmail   string              `json:"customer_notification_email"`
	DeliveryDate                DeliveryDate        `json:"delivery_date"`
	Fulfillment                 Fulfillment         `json:"fulfillment"`
	FullyRefunded               bool                `json:"fully_refunded"`
	HasCustomerMessage          bool                `json:"has_customer_message"`
	HasIncident                 bool                `json:"has_incident"`
	HasInvoice                  bool                `json:"has_invoice"`
	InvoiceDetails              InvoiceDetails      `json:"invoice_details"`
	LastUpdatedDate             string              `json:"last_updated_date"`
	LeadtimeToShip              int                 `json:"leadtime_to_ship"`
	OrderAdditionalFields       []AdditionalField   `json:"order_additional_fields"`
	OrderID                     string              `json:"order_id"`
	OrderLines                  []OrderLine         `json:"order_lines"`
	OrderState                  string              `json:"order_state"`
	OrderStateReasonCode        *string             `json:"order_state_reason_code"`
	OrderStateReasonLabel       *string             `json:"order_state_reason_label"`
	OrderTaxMode                string              `json:"order_tax_mode"`
	OrderTaxes                  []OrderTax          `json:"order_taxes"`
	PaymentType                 string              `json:"payment_type"`
	PaymentWorkflow             string              `json:"payment_workflow"`
	Price                       float64             `json:"price"`
	Promotions                  Promotions          `json:"promotions"`
	QuoteID                     *string             `json:"quote_id"`
	References                  References          `json:"references"`
	ShippingCarrierCode         string              `json:"shipping_carrier_code"`
	ShippingCarrierStandardCode string              `json:"shipping_carrier_standard_code"`
	ShippingCompany             string              `json:"shipping_company"`
	ShippingDeadline            string              `json:"shipping_deadline"`
	ShippingPrice               float64             `json:"shipping_price"`
	ShippingPudoID              string              `json:"shipping_pudo_id"`
	ShippingTracking            string              `json:"shipping_tracking"`
	ShippingTrackingURL         string              `json:"shipping_tracking_url"`
	ShippingTypeCode            string              `json:"shipping_type_code"`
	ShippingTypeLabel           string              `json:"shipping_type_label"`
	ShippingTypeStandardCode    string              `json:"shipping_type_standard_code"`
	ShippingZoneCode            string              `json:"shipping_zone_code"`
	ShippingZoneLabel           string              `json:"shipping_zone_label"`
	TotalCommission             float64             `json:"total_commission"`
	TotalPrice                  float64             `json:"total_price"`
	TransactionDate             string              `json:"transaction_date"`
	TransactionNumber           string              `json:"transaction_number"`
}

// OrderLine represents an individual line item in an order
type OrderLine struct {
	CanRefund                 bool                       `json:"can_refund"`
	Cancelations              []Cancelation              `json:"cancelations"`
	CategoryCode              string                     `json:"category_code"`
	CategoryLabel             string                     `json:"category_label"`
	CommissionFee             float64                    `json:"commission_fee"`
	CommissionTaxes           []CommissionTax            `json:"commission_taxes"`
	CreatedDate               string                     `json:"created_date"`
	DebitedDate               string                     `json:"debited_date"`
	Description               *string                    `json:"description"`
	Fees                      []Fee                      `json:"fees"`
	LastUpdatedDate           string                     `json:"last_updated_date"`
	OfferID                   int                        `json:"offer_id"`
	OfferSKU                  string                     `json:"offer_sku"`
	OfferStateCode            string                     `json:"offer_state_code"`
	OrderLineAdditionalFields []OrderLineAdditionalField `json:"order_line_additional_fields"`
	OrderLineID               string                     `json:"order_line_id"`
	OrderLineIndex            int                        `json:"order_line_index"`
	OrderLineState            string                     `json:"order_line_state"`
	OrderLineStateReasonCode  *string                    `json:"order_line_state_reason_code"`
	OrderLineStateReasonLabel *string                    `json:"order_line_state_reason_label"`
	OriginUnitPrice           float64                    `json:"origin_unit_price"`
	Price                     float64                    `json:"price"`
	PriceAdditionalInfo       *string                    `json:"price_additional_info"`
	PriceUnit                 float64                    `json:"price_unit"`
	ProductMedias             []ProductMedia             `json:"product_medias"`
	ProductShopSKU            string                     `json:"product_shop_sku"`
	ProductSKU                string                     `json:"product_sku"`
	ProductTitle              string                     `json:"product_title"`
	Promotions                []interface{}              `json:"promotions"`
	PurchaseInformation       PurchaseInformation        `json:"purchase_information"`
	Quantity                  int                        `json:"quantity"`
	ReceivedDate              string                     `json:"received_date"`
	Refunds                   []Refund                   `json:"refunds"`
	ShippedDate               string                     `json:"shipped_date"`
	ShippingFrom              ShippingFrom               `json:"shipping_from"`
	ShippingPrice             float64                    `json:"shipping_price"`
	ShippingTaxes             []ShippingTax              `json:"shipping_taxes"`
	TaxLegalNotice            string                     `json:"tax_legal_notice"`
	Taxes                     []Tax                      `json:"taxes"`
	TotalCommission           float64                    `json:"total_commission"`
	TotalPrice                float64                    `json:"total_price"`
}

// OrderLineAdditionalField represents additional information for an order line
type OrderLineAdditionalField struct {
	Code  string      `json:"code"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// Channel represents the sales channel information
type Channel struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// DeliveryDate represents delivery time frame
type DeliveryDate struct {
	Earliest string `json:"earliest"`
	Latest   string `json:"latest"`
}

// References represents order references
type References struct {
	OrderReferenceForCustomer string `json:"order_reference_for_customer"`
	OrderReferenceForSeller   string `json:"order_reference_for_seller"`
}

// Promotions represents promotion information
type Promotions struct {
	AppliedPromotions   []interface{} `json:"applied_promotions"`
	TotalDeducedAmount  float64       `json:"total_deduced_amount"`
}

// Orders API client
type Orders struct {
	Client RequestHandler
}

// List retrieves all orders
func (o Orders) List(params url.Values) ([]Order, error) {
	type response struct {
		Orders     []Order `json:"orders"`
		TotalCount int     `json:"total_count"`
	}
	
	var resp response
	err := o.Client.Request("GET", "orders", params, &resp)
	return resp.Orders, err
}

// GetOrderLines retrieves order lines for a given order
func (o Orders) GetOrderLines(id string) ([]OrderLine, error) {
	order, err := o.Get(id)
	if err != nil {
		return nil, err
	}
	return order.OrderLines, nil
}

// Get retrieves a single order by ID
func (o Orders) Get(id string) (Order, error) {
	type response struct {
		Orders []Order `json:"orders"`
	}
	
	var resp response
	var order Order
	err := o.Client.Request("GET", "orders/"+id, nil, &resp)
	if err == nil && len(resp.Orders) > 0 {
		order = resp.Orders[0]
	}
	return order, err
}
