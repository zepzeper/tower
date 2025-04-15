package woocommerce

import (
	"net/url"
	"testing"

	"github.com/zepzeper/tower/internal/config"
)

func getTestClient() Client {
	cc := ClientConfig{
		APIHost:        config.GetEnv("WC_API_HOST", ""),
		ConsumerKey:    config.GetEnv("WC_PB_KEY", ""),
		ConsumerSecret: config.GetEnv("WC_PP_KEY", ""),
		Debug:          false,
	}
	return NewClient(cc)
}

func TestClientConfig(t *testing.T) {
	check := "https://example.com/wp-json/wc/v3/products"
	cc := ClientConfig{
		APIHost: "https://example.com/wp-json/wc/v3",
	}

	r := cc.GetAPIEndpoint("products")
	if r != check {
		t.Errorf("Expected '%s', got '%s'", check, r)
	}
}

func TestGetCustomers(t *testing.T) {
	c := getTestClient()

	params := url.Values{}
	params.Set("orderby", "id")
	params.Set("order", "desc")

	r, err := c.Customers().List(params)
	if err != nil {
		t.Fatal(err)
	}

	if len(r) > 0 {
		_, err = c.Customers().Get(r[0].ID)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetOrders(t *testing.T) {
	c := getTestClient()

	params := url.Values{}
	params.Set("orderby", "id")
	params.Set("order", "desc")

	r, err := c.Orders().List(params)
	if err != nil {
		t.Fatal(err)
	}

	if len(r) > 0 {
		_, err = c.Orders().Get(r[0].ID)
		if err != nil {
			t.Error(err)
		}

		ons, err := c.Orders().ListOrderNotes(r[0].ID)
		if err != nil {
			t.Error(err)
		}
		if len(ons) > 0 {
			_, err = c.Orders().GetOrderNote(r[0].ID, ons[0].ID)

		  if err != nil {
		  	t.Error(err)
	  	}
		}
	}
}
