package mirakl

import (
	"net/url"
	"testing"

	"github.com/zepzeper/tower/internal/config"
)

func getTestClient() Client {
	cc := ClientConfig{
		Marketplace:    config.GetEnv("MK_MP", ""),
		ShopID:         config.GetEnv("MK_SHOP_KEY", ""),
    ApiKey:         config.GetEnv("MK_PP_KEY", ""),
		Debug:          true,
	}
	return NewClient(cc)
}

func TestGetProducts(t *testing.T) {
	c := getTestClient()

	params := url.Values{}

	_, err := c.Products().List(params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOrders(t *testing.T) {
	c := getTestClient()

	params := url.Values{}

	_, err := c.Orders().List(params)
	if err != nil {
		t.Fatal(err)
	}
}
