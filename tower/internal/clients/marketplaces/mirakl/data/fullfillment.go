package data

// Fulfillment represents fulfillment information
type Fulfillment struct {
	Center FulfillmentCenter `json:"center"`
}

// FulfillmentCenter represents a fulfillment center
type FulfillmentCenter struct {
	Code string `json:"code"`
}

// ShippingFrom represents shipping origin information
type ShippingFrom struct {
	Address   Address  `json:"address"`
	Warehouse *string  `json:"warehouse"`
}
