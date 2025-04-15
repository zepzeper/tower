package data

// Customer represents customer information
type Customer struct {
	AccountingContact Contact      `json:"accounting_contact"`
	BillingAddress    Address      `json:"billing_address"`
	Civility          string       `json:"civility"`
	CustomerID        string       `json:"customer_id"`
	DeliveryContact   Contact      `json:"delivery_contact"`
	Firstname         string       `json:"firstname"`
	Lastname          string       `json:"lastname"`
	Locale            string       `json:"locale"`
	Organization      Organization `json:"organization"`
	ShippingAddress   Address      `json:"shipping_address"`
}

// Contact represents a contact person
type Contact struct {
	CustomerID string `json:"customer_id"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Locale     string `json:"locale"`
}

// Address represents a physical address
type Address struct {
	AdditionalInfo  string `json:"additional_info,omitempty"`
	City            string `json:"city"`
	Civility        string `json:"civility,omitempty"`
	Company         string `json:"company,omitempty"`
	Company2        string `json:"company_2,omitempty"`
	Country         string `json:"country,omitempty"`
	CountryIsoCode  string `json:"country_iso_code"`
	Firstname       string `json:"firstname,omitempty"`
	Lastname        string `json:"lastname,omitempty"`
	State           string `json:"state"`
	Street1         string `json:"street_1"`
	Street2         string `json:"street_2,omitempty"`
	ZipCode         string `json:"zip_code"`
}

// Organization represents a business organization
type Organization struct {
	Address                 Address `json:"address"`
	IdentificationNumber    string  `json:"identification_number"`
	Name                    string  `json:"name"`
	OrganizationID          string  `json:"organization_id"`
	TaxIdentificationNumber string  `json:"tax_identification_number"`
}
