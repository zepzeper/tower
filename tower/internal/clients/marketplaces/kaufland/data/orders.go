package data

import (
	"net/url"
	"strconv"
	"time"
)

type Order struct {
    AgeRating                   int             `json:"age_rating"`
    BillingAddress              BillingAddress  `json:"billing_address"`
    Buyer                       Buyer           `json:"buyer"`
    CancelReason                any             `json:"cancel_reason"`
    DangerLabel9a               any             `json:"danger_label_9A"`
    DangerousGoodsLiShipping    any             `json:"dangerous_goods_li_shipping"`
    DeliveryTimeExpiresIso      time.Time       `json:"delivery_time_expires_iso,omitzero"`
    DeliveryTimeMax             int             `json:"delivery_time_max,omitzero"`
    DeliveryTimeMin             int             `json:"delivery_time_min,omitzero"`
    IDCategory                  int             `json:"id_category,omitzero"`
    IDOffer                     string          `json:"id_offer,omitzero"`
    IDOrder                     string          `json:"id_order,omitzero"`
    IDOrderUnit                 int             `json:"id_order_unit,omitzero"`
    IsMarketplaceDeemedSupplier bool            `json:"is_marketplace_deemed_supplier"`
    IsValid                     bool            `json:"is_valid,omitzero"`
    MainPicture                 string          `json:"main_picture,omitzero"`
    Manufacturer                string          `json:"manufacturer,omitzero"`
    Note                        string          `json:"note,omitzero"`
    Price                       int             `json:"price,omitzero"`
    Product                     Product         `json:"product"`
    RealMgbArticleNumber        string          `json:"real_mgb_article_number"`
    RevenueGross                int             `json:"revenue_gross,omitzero"`
    RevenueNet                  int             `json:"revenue_net,omitzero"`
    ShippingAddress             ShippingAddress `json:"shipping_address"`
    ShippingRate                int             `json:"shipping_rate"`
    Status                      string          `json:"status,omitzero"`
    Storefront                  string          `json:"storefront,omitzero"`
    TsCreatedIso                time.Time       `json:"ts_created_iso,omitzero"`
    TsUpdatedIso                time.Time       `json:"ts_updated_iso,omitzero"`
    UnitCondition               string          `json:"unit_condition,omitzero"`
    URL                         string          `json:"url,omitzero"`
    Vat                         int             `json:"vat,omitzero"`

}

type ShippingAddress struct {
    AdditionalField string `json:"additional_field,omitzero"`
    City            string `json:"city,omitzero"`
    CompanyName     string `json:"company_name"`
    Country         string `json:"country,omitzero"`
    FirstName       string `json:"first_name,omitzero"`
    HouseNumber     string `json:"house_number,omitzero"`
    LastName        string `json:"last_name,omitzero"`
    Phone           string `json:"phone,omitzero"`
    Postcode        string `json:"postcode,omitzero"`
    Street          string `json:"street,omitzero"`
} 

type BillingAddress struct {
    AdditionalField string `json:"additional_field,omitzero"`
    City            string `json:"city,omitzero"`
    CompanyName     string `json:"company_name"`
    Country         string `json:"country,omitzero"`
    FirstName       string `json:"first_name,omitzero"`
    HouseNumber     string `json:"house_number,omitzero"`
    LastName        string `json:"last_name,omitzero"`
    Phone           string `json:"phone,omitzero"`
    Postcode        string `json:"postcode,omitzero"`
    Street          string `json:"street,omitzero"`
} 

type Buyer struct {
    Email   string `json:"email,omitzero"`
    IDBuyer int    `json:"id_buyer,omitzero"`
} 


type Orders struct {
	Client RequestHandler
}


// List all orders
func (o Orders) List(listParams url.Values) ([]Order, error) {
	r := make([]Order, 0)
	err := o.Client.Request("GET", "order-units", listParams, &r)
	return r, err
}

// Get order by ID
func (o Orders) Get(id int) (Order, error) {
	var r Order
	err := o.Client.Request("GET", "orders-units/"+strconv.Itoa(id), nil, &r)
	return r, err
}
