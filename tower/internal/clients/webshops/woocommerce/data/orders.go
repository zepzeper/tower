package data

import (
	"net/url"
	"strconv"
)

// Order data
type Order struct {
	ID                 int                 `json:"id"`
	ParentID           int                 `json:"parent_id"`
	Number             string              `json:"number"`
	OrderKey           string              `json:"order_key"`
	CreatedVia         string              `json:"created_via"`
	Version            string              `json:"version"`
	Status             string              `json:"status"`
	Currency           string              `json:"currency"`
	DateCreated        string              `json:"date_created"`
	DateCreatedGmt     string              `json:"date_created_gmt"`
	DateModified       string              `json:"date_modified"`
	DateModifiedGmt    string              `json:"date_modified_gmt"`
	DiscountTotal      string              `json:"discount_total"`
	DiscountTax        string              `json:"discount_tax"`
	ShippingTotal      string              `json:"shipping_total"`
	ShippingTax        string              `json:"shipping_tax"`
	CartTax            string              `json:"cart_tax"`
	Total              string              `json:"total"`
	TotalTax           string              `json:"total_tax"`
	PricesIncludeTax   bool                `json:"prices_include_tax"`
	CustomerID         int                 `json:"customer_id"`
	CustomerIPAddress  string              `json:"customer_ip_address"`
	CustomerUserAgent  string              `json:"customer_user_agent"`
	CustomerNote       string              `json:"customer_note"`
	Billing            Billing             `json:"billing"`
	Shipping           Shipping            `json:"shipping"`
	PaymentMethod      string              `json:"payment_method"`
	PaymentMethodTitle string              `json:"payment_method_title"`
	TransactionID      string              `json:"transaction_id"`
	DatePaid           string              `json:"date_paid"`
	DatePaidGmt        string              `json:"date_paid_gmt"`
	DateCompleted      string              `json:"date_completed"`
	DateCompletedGmt   string              `json:"date_completed_gmt"`
	CartHash           string              `json:"cart_hash"`
	MetaData           []MetaData          `json:"meta_data"`
	LineItems          []OrderLineItem     `json:"line_items"`
	TaxLines           []OrderTaxLine      `json:"tax_lines"`
	ShippingLines      []OrderShippingLine `json:"shipping_lines"`
	FeeLines           []OrderFeeLine      `json:"fee_lines"`
	CouponLines        []OrderCouponLine   `json:"coupon_lines"`
	Refunds            []OrderRefundLine   `json:"refunds"`
	Links              Links               `json:"_links"`
}

// OrderLineItem data
type OrderLineItem struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	ProductID   int                `json:"product_id"`
	VariationID int                `json:"variation_id"`
	Quantity    int                `json:"quantity"`
	TaxClass    string             `json:"tax_class"`
	Subtotal    string             `json:"subtotal"`
	SubtotalTax string             `json:"subtotal_tax"`
	Total       string             `json:"total"`
	TotalTax    string             `json:"total_tax"`
	Taxes       []OrderLineItemTax `json:"taxes"`
	MetaData    []MetaData         `json:"meta_data"`
	Sku         string             `json:"sku"`
	Price       float32            `json:"price"`
}

// OrderLineItemTax data
type OrderLineItemTax struct {
	ID       int    `json:"id"`
	Total    float32 `json:"total"`
	Subtotal float32 `json:"subtotal"`
}

// OrderTaxLine data
type OrderTaxLine struct {
	ID               int        `json:"id"`
	RateCode         string     `json:"rate_code"`
	RateID           int        `json:"rate_id"`
	Label            string     `json:"label"`
	Compound         bool       `json:"compound"`
	TaxTotal         float32    `json:"tax_total"`
	ShippingTaxTotal float32    `json:"shipping_tax_total"`
	MetaData         []MetaData `json:"meta_data"`
}

// OrderShippingLine data
type OrderShippingLine struct {
	ID          int            `json:"id"`
	MethodTitle string         `json:"method_title"`
	MethodID    string         `json:"method_id"`
	Total       float32        `json:"total"`
	TotalTax    float32        `json:"total_tax"`
	Taxes       []OrderTaxLine `json:"taxes"`
	MetaData    []MetaData     `json:"meta_data"`
}

// OrderFeeLine data
type OrderFeeLine struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	TaxClass  string         `json:"tax_class"`
	TaxStatus string         `json:"tax_status"`
	Total     float32        `json:"total"`
	TotalTax  float32        `json:"total_tax"`
	Taxes     []OrderTaxLine `json:"taxes"`
	MetaData  []MetaData     `json:"meta_data"`
}

// OrderCouponLine data
type OrderCouponLine struct {
	ID          int        `json:"id"`
	Code        string     `json:"code"`
	Discount    string     `json:"discount"`
	DiscountTax string     `json:"discount_tax"`
	MetaData    []MetaData `json:"meta_data"`
}

// OrderRefundLine data
type OrderRefundLine struct {
	ID     int     `json:"id"`
	Reason string  `json:"reason"`
	Total  float32 `json:"total"`
}

// OrderNote data
type OrderNote struct {
	ID             int    `json:"id"`
	Author         string `json:"author"`
	DateCreated    string `json:"date_created"`
	DateCreatedGmt string `json:"date_created_gmt"`
	Note           string `json:"note"`
	CustomerNote   bool   `json:"customer_note"`
	Links          Links  `json:"_links"`
}

// Orders API client
type Orders struct {
    Client RequestHandler
}

// List all orders
func (o Orders) List(listParams url.Values) ([]Order, error) {
	r := make([]Order, 0)
	err := o.Client.Request("GET", "orders", listParams, &r)
	return r, err
}

// Get order by ID
func (o Orders) Get(id int) (Order, error) {
	var r Order
	err := o.Client.Request("GET", "orders/"+strconv.Itoa(id), nil, &r)
	return r, err
}

// ListOrderNotes retrieves all notes for a given order
func (o Orders) ListOrderNotes(id int) ([]OrderNote, error) {
	r := make([]OrderNote, 0)
	err := o.Client.Request("GET", "orders/"+strconv.Itoa(id)+"/notes", nil, &r)
	return r, err
}

// GetOrderNote by IDs
func (o Orders) GetOrderNote(orderID, noteID int) (OrderNote, error) {
	var r OrderNote
	err := o.Client.Request("GET", "orders/"+strconv.Itoa(orderID)+"/notes/"+strconv.Itoa(noteID), nil, &r)
	return r, err
}
