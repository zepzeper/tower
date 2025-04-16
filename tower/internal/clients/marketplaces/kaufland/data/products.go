package data

import (
	"net/url"
	"strconv"
	"time"
)

type Product struct {
        Amount            int       `json:"amount,omitzero"`
        Condition         string    `json:"condition,omitzero"`
        Currency          string    `json:"currency,omitzero"`
        DateInsertedIso   time.Time `json:"date_inserted_iso,omitzero"`
        DateLastchangeIso time.Time `json:"date_lastchange_iso,omitzero"`
        FulfillmentType   string    `json:"fulfillment_type,omitzero"`
        HandlingTime      int       `json:"handling_time,omitzero"`
        IDOffer           string    `json:"id_offer,omitzero"`
        IDProduct         int       `json:"id_product,omitzero"`
        IDShippingGroup   int       `json:"id_shipping_group,omitzero"`
        IDUnit            int       `json:"id_unit,omitzero"`
        IDWarehouse       int       `json:"id_warehouse,omitzero"`
        ListingPrice      int       `json:"listing_price,omitzero"`
        MinimumPrice      int       `json:"minimum_price,omitzero"`
        Note              string    `json:"note,omitzero"`
        Price             int       `json:"price,omitzero"`
        Product           ProductInformation `json:"product,omitzero"`
        ShippingRate     int    `json:"shipping_rate"`
        Status           string `json:"status,omitzero"`
        Storefront       string `json:"storefront,omitzero"`
        TransportTimeMax int    `json:"transport_time_max,omitzero"`
        TransportTimeMin int    `json:"transport_time_min,omitzero"`
        VatIndicator     string `json:"vat_indicator,omitzero"`
}

type ProductInformation struct {
        AgeRating                int      `json:"age_rating"`
        DangerLabel9a            string   `json:"danger_label_9A,omitzero"`
        DangerousGoodsLiShipping string   `json:"dangerous_goods_li_shipping,omitzero"`
        Eans                     []string `json:"eans"`
        IDCategory               int      `json:"id_category,omitzero"`
        IDProduct                int      `json:"id_product,omitzero"`
        IsValid                  bool     `json:"is_valid,omitzero"`
        MainPicture              string   `json:"main_picture,omitzero"`
        Manufacturer             string   `json:"manufacturer,omitzero"`
        Storefront               string   `json:"storefront,omitzero"`
        Title                    string   `json:"title,omitzero"`
        URL                      string   `json:"url,omitzero"`
}


type Products struct {
	Client RequestHandler
}


func (c Products) List(listParams url.Values) ([]Product, error) {
	r := make([]Product, 0)
	err := c.Client.Request("GET", "units", listParams, &r)
	return r, err
}

// Get customer by ID
func (c Products) Get(id int) (Product, error) {
	var or Product
	err := c.Client.Request("GET", "units/"+strconv.Itoa(id), nil, &or)
	return or, err
}
