package woocommerce

import "github.com/zepzeper/tower/internal/core/connectors"

// GetSchema returns the schema for WooCommerce
func GetSchema() connectors.Schema {
	return connectors.Schema{
		EntityName: "product",
		Fields: map[string]connectors.FieldDefinition{
			"id": {
				Type:     "number",
				Required: false,
				Path:     "id",
			},
			"name": {
				Type:     "string",
				Required: true,
				Path:     "name",
			},
			"slug": {
				Type:     "string",
				Required: false,
				Path:     "slug",
			},
			"permalink": {
				Type:     "string",
				Required: false,
				Path:     "permalink",
			},
			"date_created": {
				Type:     "string",
				Required: false,
				Path:     "date_created",
			},
			"date_modified": {
				Type:     "string",
				Required: false,
				Path:     "date_modified",
			},
			"type": {
				Type:     "string",
				Required: true,
				Path:     "type",
				EnumValues: []string{"simple", "grouped", "external", "variable"},
			},
			"status": {
				Type:     "string",
				Required: false,
				Path:     "status",
				EnumValues: []string{"draft", "pending", "private", "publish"},
			},
			"featured": {
				Type:     "boolean",
				Required: false,
				Path:     "featured",
			},
			"catalog_visibility": {
				Type:     "string",
				Required: false,
				Path:     "catalog_visibility",
			},
			"description": {
				Type:     "string",
				Required: false,
				Path:     "description",
			},
			"short_description": {
				Type:     "string",
				Required: false,
				Path:     "short_description",
			},
			"sku": {
				Type:     "string",
				Required: false,
				Path:     "sku",
			},
			"price": {
				Type:     "number",
				Required: true,
				Path:     "price",
			},
			"regular_price": {
				Type:     "number",
				Required: false,
				Path:     "regular_price",
			},
			"sale_price": {
				Type:     "number",
				Required: false,
				Path:     "sale_price",
			},
			"on_sale": {
				Type:     "boolean",
				Required: false,
				Path:     "on_sale",
			},
			"purchasable": {
				Type:     "boolean",
				Required: false,
				Path:     "purchasable",
			},
			"total_sales": {
				Type:     "number",
				Required: false,
				Path:     "total_sales",
			},
			"virtual": {
				Type:     "boolean",
				Required: false,
				Path:     "virtual",
			},
			"downloadable": {
				Type:     "boolean",
				Required: false,
				Path:     "downloadable",
			},
			"downloads": {
				Type:     "array",
				Required: false,
				Path:     "downloads",
			},
			"download_limit": {
				Type:     "number",
				Required: false,
				Path:     "download_limit",
			},
			"download_expiry": {
				Type:     "number",
				Required: false,
				Path:     "download_expiry",
			},
			"external_url": {
				Type:     "string",
				Required: false,
				Path:     "external_url",
			},
			"button_text": {
				Type:     "string",
				Required: false,
				Path:     "button_text",
			},
			"tax_status": {
				Type:     "string",
				Required: false,
				Path:     "tax_status",
			},
			"tax_class": {
				Type:     "string",
				Required: false,
				Path:     "tax_class",
			},
			"manage_stock": {
				Type:     "boolean",
				Required: false,
				Path:     "manage_stock",
			},
			"stock_quantity": {
				Type:     "number",
				Required: false,
				Path:     "stock_quantity",
			},
			"stock_status": {
				Type:     "string",
				Required: false,
				Path:     "stock_status",
			},
			"backorders": {
				Type:     "string",
				Required: false,
				Path:     "backorders",
			},
			"backorders_allowed": {
				Type:     "boolean",
				Required: false,
				Path:     "backorders_allowed",
			},
			"backordered": {
				Type:     "boolean",
				Required: false,
				Path:     "backordered",
			},
			"sold_individually": {
				Type:     "boolean",
				Required: false,
				Path:     "sold_individually",
			},
			"weight": {
				Type:     "string",
				Required: false,
				Path:     "weight",
			},
			"dimensions": {
				Type:     "object",
				Required: false,
				Path:     "dimensions",
			},
			"shipping_required": {
				Type:     "boolean",
				Required: false,
				Path:     "shipping_required",
			},
			"shipping_taxable": {
				Type:     "boolean",
				Required: false,
				Path:     "shipping_taxable",
			},
			"shipping_class": {
				Type:     "string",
				Required: false,
				Path:     "shipping_class",
			},
			"shipping_class_id": {
				Type:     "number",
				Required: false,
				Path:     "shipping_class_id",
			},
			"reviews_allowed": {
				Type:     "boolean",
				Required: false,
				Path:     "reviews_allowed",
			},
			"average_rating": {
				Type:     "string",
				Required: false,
				Path:     "average_rating",
			},
			"rating_count": {
				Type:     "number",
				Required: false,
				Path:     "rating_count",
			},
			"related_ids": {
				Type:     "array",
				Required: false,
				Path:     "related_ids",
			},
			"upsell_ids": {
				Type:     "array",
				Required: false,
				Path:     "upsell_ids",
			},
			"cross_sell_ids": {
				Type:     "array",
				Required: false,
				Path:     "cross_sell_ids",
			},
			"parent_id": {
				Type:     "number",
				Required: false,
				Path:     "parent_id",
			},
			"purchase_note": {
				Type:     "string",
				Required: false,
				Path:     "purchase_note",
			},
			"categories": {
				Type:     "array",
				Required: false,
				Path:     "categories",
			},
			"tags": {
				Type:     "array",
				Required: false,
				Path:     "tags",
			},
			"images": {
				Type:     "array",
				Required: false,
				Path:     "images",
			},
			"attributes": {
				Type:     "array",
				Required: false,
				Path:     "attributes",
			},
			"default_attributes": {
				Type:     "array",
				Required: false,
				Path:     "default_attributes",
			},
			"variations": {
				Type:     "array",
				Required: false,
				Path:     "variations",
			},
			"grouped_products": {
				Type:     "array",
				Required: false,
				Path:     "grouped_products",
			},
			"menu_order": {
				Type:     "number",
				Required: false,
				Path:     "menu_order",
			},
		},
	}
}
