package cdm

import (
	"time"
)

// Entity represents a base entity type in the canonical data model
type Entity struct {
	ID        string                 `json:"id"`
	EntityType string                `json:"entityType"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Product represents a product in the canonical data model
type Product struct {
	Entity
	SKU              string   `json:"sku"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	Price            float64  `json:"price"`
	Currency         string   `json:"currency,omitempty"`
	Images           []string `json:"images,omitempty"`
	Categories       []string `json:"categories,omitempty"`
	Status           string   `json:"status"`
	StockQuantity    int      `json:"stockQuantity"`
	Weight           float64  `json:"weight,omitempty"`
	WeightUnit       string   `json:"weightUnit,omitempty"`
	Dimensions       *Dimensions `json:"dimensions,omitempty"`
	Attributes       map[string]interface{} `json:"attributes,omitempty"`
	Variations       []string `json:"variations,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

// Dimensions represents product dimensions
type Dimensions struct {
	Length float64 `json:"length,omitempty"`
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
	Unit   string  `json:"unit,omitempty"`
}

// Customer represents a customer in the canonical data model
type Customer struct {
	Entity
	Email         string    `json:"email"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Company       string    `json:"company,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	BillingAddress *Address `json:"billingAddress,omitempty"`
	ShippingAddress *Address `json:"shippingAddress,omitempty"`
	Notes         string    `json:"notes,omitempty"`
	IsActive      bool      `json:"isActive"`
}

// Address represents a physical address
type Address struct {
	Line1       string `json:"line1"`
	Line2       string `json:"line2,omitempty"`
	City        string `json:"city"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postalCode"`
	Country     string `json:"country"`
}

// Order represents an order in the canonical data model
type Order struct {
	Entity
	OrderNumber    string       `json:"orderNumber"`
	CustomerID     string       `json:"customerId"`
	OrderStatus    string       `json:"orderStatus"`
	Currency       string       `json:"currency"`
	Subtotal       float64      `json:"subtotal"`
	ShippingAmount float64      `json:"shippingAmount,omitempty"`
	TaxAmount      float64      `json:"taxAmount,omitempty"`
	DiscountAmount float64      `json:"discountAmount,omitempty"`
	TotalAmount    float64      `json:"totalAmount"`
	PaymentMethod  string       `json:"paymentMethod,omitempty"`
	ShippingMethod string       `json:"shippingMethod,omitempty"`
	LineItems      []OrderItem  `json:"lineItems"`
	BillingAddress *Address     `json:"billingAddress,omitempty"`
	ShippingAddress *Address    `json:"shippingAddress,omitempty"`
	Notes          string       `json:"notes,omitempty"`
	PlacedDate     time.Time    `json:"placedDate"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID      string  `json:"productId"`
	SKU            string  `json:"sku"`
	Name           string  `json:"name"`
	Quantity       int     `json:"quantity"`
	UnitPrice      float64 `json:"unitPrice"`
	Subtotal       float64 `json:"subtotal"`
	TaxAmount      float64 `json:"taxAmount,omitempty"`
	DiscountAmount float64 `json:"discountAmount,omitempty"`
	Total          float64 `json:"total"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// Convert converts a generic data payload to a specific entity type
func Convert(entityType string, data map[string]interface{}) (interface{}, error) {
	switch entityType {
	case "product":
		// Convert map to Product
		return ConvertToProduct(data), nil
	case "customer":
		// Convert map to Customer
		return ConvertToCustomer(data), nil
	case "order":
		// Convert map to Order
		return ConvertToOrder(data), nil
	default:
		// Return generic map as is
		return data, nil
	}
}

// ConvertFromEntity converts an entity to a generic data payload
func ConvertFromEntity(entity interface{}) map[string]interface{} {
	// Implementation will depend on your serialization library
	// This is a placeholder
	switch e := entity.(type) {
	case Product:
		// Convert Product to map
		return map[string]interface{}{
			"id":           e.ID,
			"entityType":   e.EntityType,
			"sku":          e.SKU,
			"name":         e.Name,
			// ... other fields
		}
	case Customer:
		// Convert Customer to map
		return map[string]interface{}{
			"id":           e.ID,
			"entityType":   e.EntityType,
			"email":        e.Email,
			// ... other fields
		}
	// ... other entity types
	default:
		// Return empty map if not a known entity
		return map[string]interface{}{}
	}
}

// Helper functions to convert maps to entity types
func ConvertToProduct(data map[string]interface{}) Product {
	// Implementation will depend on your serialization library
	// This is a placeholder implementation
	product := Product{}
	// Populate fields from map
	if id, ok := data["id"].(string); ok {
		product.ID = id
	}
	if sku, ok := data["sku"].(string); ok {
		product.SKU = sku
	}
	// ... and so on for other fields
	
	product.EntityType = "product"
	return product
}

func ConvertToCustomer(data map[string]interface{}) Customer {
	// Implementation will depend on your serialization library
	// This is a placeholder implementation
	customer := Customer{}
	// Populate fields from map
	if id, ok := data["id"].(string); ok {
		customer.ID = id
	}
	if email, ok := data["email"].(string); ok {
		customer.Email = email
	}
	// ... and so on for other fields
	
	customer.EntityType = "customer"
	return customer
}

func ConvertToOrder(data map[string]interface{}) Order {
	// Implementation will depend on your serialization library
	// This is a placeholder implementation
	order := Order{}
	// Populate fields from map
	if id, ok := data["id"].(string); ok {
		order.ID = id
	}
	if number, ok := data["orderNumber"].(string); ok {
		order.OrderNumber = number
	}
	// ... and so on for other fields
	
	order.EntityType = "order"
	return order
}
