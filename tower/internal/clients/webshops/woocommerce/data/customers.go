package data

import (
	"net/url"
	"strconv"
)

// Customer data
type Customer struct {
	ID               int        `json:"id"`
	DateCreated      string     `json:"date_created"`
	DateCreatedGmt   string     `json:"date_created_gmt"`
	DateModified     string     `json:"date_modified"`
	DateModifiedGmt  string     `json:"date_modified_gmt"`
	Email            string     `json:"email"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Role             string     `json:"role"`
	Username         string     `json:"username"`
	Billing          Billing    `json:"billing"`
	Shipping         Shipping   `json:"shipping"`
	IsPayingCustomer bool       `json:"is_paying_customer"`
	AvatarURL        string     `json:"avatar_url"`
	MetaData         []MetaData `json:"meta_data"`
	Links            Links      `json:"_links"`
}

// Customers API client
type Customers struct {
  Client RequestHandler
}


// List all customers
func (c Customers) List(listParams url.Values) ([]Customer, error) {
	r := make([]Customer, 0)
	err := c.Client.Request("GET", "customers", listParams, &r)
	return r, err
}

// Get customer by ID
func (c Customers) Get(id int) (Customer, error) {
	var or Customer
	err := c.Client.Request("GET", "customers/"+strconv.Itoa(id), nil, &or)
	return or, err
}
