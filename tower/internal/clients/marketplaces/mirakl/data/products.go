package data

import "net/url"

type Product struct {
	Active             bool                  `json:"active"`
	AllPrices          []Price               `json:"all_prices"`
	AllowQuoteRequests bool                  `json:"allow_quote_requests"`
	ApplicablePricing  Price                 `json:"applicable_pricing"`
	AvailableEndDate   string                `json:"available_end_date"`
	AvailableStartDate string                `json:"available_start_date"`
	CategoryCode       string                `json:"category_code"`
	CategoryLabel      string                `json:"category_label"`
	Channels           []string              `json:"channels"`
	CurrencyIsoCode    string                `json:"currency_iso_code"`
	Description        string                `json:"description"`
	Discount           Discount              `json:"discount"`
	Fulfillment        Fulfillment           `json:"fulfillment"`
	InactivityReasons  []string              `json:"inactivity_reasons"`
	LeadtimeToShip     int                   `json:"leadtime_to_ship"`
	LogisticClass      LogisticClass         `json:"logistic_class"`
	MinOrderQuantity   int                   `json:"min_order_quantity"`
	MinShippingPrice   float64               `json:"min_shipping_price"`
	MinShippingPriceAdditional float64       `json:"min_shipping_price_additional"`
	MinShippingType    string                `json:"min_shipping_type"`
	MinShippingZone    string                `json:"min_shipping_zone"`
	OfferAdditionalFields []AdditionalField  `json:"offer_additional_fields"`
	OfferID            int                   `json:"offer_id"`
	PackageQuantity    int                   `json:"package_quantity"`
	Price              float64               `json:"price"`
	PriceAdditionalInfo string               `json:"price_additional_info"`
	ProductBrand       string                `json:"product_brand"`
	ProductDescription string                `json:"product_description"`
	ProductReferences  []ProductReference    `json:"product_references"`
	ProductSKU         string                `json:"product_sku"`
	ProductTitle       string                `json:"product_title"`
	Quantity           int                   `json:"quantity"`
	ShippingDeadline   string                `json:"shipping_deadline"`
	ShopSKU            string                `json:"shop_sku"`
	StateCode          string                `json:"state_code"`
	TotalPrice         float64               `json:"total_price"`
}

type ProductsResponse struct {
	Products   []Product `json:"offers"`
	TotalCount int       `json:"total_count"`
}

// ProductMedia represents product media information
type ProductMedia struct {
	MediaURL string `json:"media_url"`
	MimeType string `json:"mime_type"`
	Type     string `json:"type"`
}

// Price represents pricing information with optional context
type Price struct {
	ChannelCode        string                `json:"channel_code,omitempty"`
	Context            *PriceContext         `json:"context,omitempty"`
	DiscountEndDate    string                `json:"discount_end_date"`
	DiscountStartDate  string                `json:"discount_start_date"`
	Price              float64               `json:"price"`
	UnitDiscountPrice  float64               `json:"unit_discount_price"`
	UnitOriginPrice    float64               `json:"unit_origin_price"`
	VolumePrices       []VolumePrice         `json:"volume_prices"`
}

// PriceContext represents the context for a price application
type PriceContext struct {
	ChannelCodes           []string `json:"channel_codes,omitempty"`
	CustomerGroupIDs       []string `json:"customer_group_ids,omitempty"`
	CustomerOrganizationIDs []string `json:"customer_organization_ids,omitempty"`
	EndDate                string   `json:"end_date,omitempty"`
	StartDate              string   `json:"start_date,omitempty"`
}

// VolumePrice represents price tiers based on quantity
type VolumePrice struct {
	Price              float64 `json:"price"`
	QuantityThreshold  int     `json:"quantity_threshold"`
	UnitDiscountPrice  float64 `json:"unit_discount_price"`
	UnitOriginPrice    float64 `json:"unit_origin_price"`
}

// Discount represents discount information
type Discount struct {
	DiscountPrice      *float64      `json:"discount_price"`
	EndDate            string        `json:"end_date"`
	OriginPrice        float64       `json:"origin_price"`
	Ranges             []PriceRange  `json:"ranges"`
	StartDate          string        `json:"start_date"`
}

// PriceRange represents a price tier in a discount
type PriceRange struct {
	Price              float64 `json:"price"`
	QuantityThreshold  int     `json:"quantity_threshold"`
}

// LogisticClass represents logistic classification
type LogisticClass struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// AdditionalField represents additional information for an offer
type AdditionalField struct {
	Code  string `json:"code"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ProductReference represents a product identifier
type ProductReference struct {
	Reference     string `json:"reference"`
	ReferenceType string `json:"reference_type"`
}

type Products struct {
	Client RequestHandler
}

// List retrieves all offers
func (p Products) List(params url.Values) ([]Product, error) {
	var resp ProductsResponse
	err := p.Client.Request("GET", "/offers", params, &resp)
	return resp.Products, err
}

// Get retrieves a single product by ID
func (p Products) Get(id int64) (Product, error) {
	var product Product
	err := p.Client.Request("GET", "/offers/"+string(id), nil, &product)
	return product, err
}
