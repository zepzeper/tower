package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zepzeper/tower/internal/connectors/client"
	"github.com/zepzeper/tower/internal/helpers"
)

func main() {
	helpers.LoadEnvFile()
	// Check environment variables
	checkRequiredEnvVars()
	
	// Initialize client factory
	clientFactory := client.NewFactory()
	
	// Step 1: Create WooCommerce client
	log.Println("Creating WooCommerce client...")
	wooClient, err := clientFactory.CreateClient("woocommerce", false)
	if err != nil {
		log.Fatalf("Failed to create WooCommerce client: %v", err)
	}
	
	// Step 2: Test WooCommerce connection
	log.Println("Testing WooCommerce connection...")
	_, err = wooClient.TestRequest()
	if err != nil {
		log.Fatalf("WooCommerce connection test failed: %v", err)
	}
	log.Println("WooCommerce connection successful!")
	
	// Step 3: Create Brincr client
	log.Println("Creating Brincr client...")
	brincrClient, err := clientFactory.CreateClient("brincr", true)
	if err != nil {
		log.Fatalf("Failed to create Brincr client: %v", err)
	}
	
	// Step 4: Test Brincr connection
	log.Println("Testing Brincr connection...")
	_, err = brincrClient.TestRequest()
	if err != nil {
		log.Fatalf("Brincr connection test failed: %v", err)
	}
	log.Println("Brincr connection successful!")
	
	// Step 5: Fetch products from WooCommerce
	log.Println("Fetching products from WooCommerce...")
	var wooProducts []map[string]interface{}
	
	err = wooClient.Execute("GET", "/products?per_page=5&status=publish", nil, &wooProducts)
	if err != nil {
		log.Fatalf("Failed to fetch products from WooCommerce: %v", err)
	}
	
	log.Printf("Successfully fetched %d products from WooCommerce", len(wooProducts))
	
	if len(wooProducts) == 0 {
		log.Println("No products found in WooCommerce. Exiting.")
		return
	}
	
	// Print first product for verification
	prettyPrint("WooCommerce Product Example:", wooProducts[0])
	
	// Step 6: Transform products for Brincr
	log.Println("Transforming products for Brincr...")
	brincrProducts := make([]map[string]interface{}, 0, len(wooProducts))
	
	for _, wooProduct := range wooProducts {
		// Create a basic Brincr product
		brincrProduct := map[string]interface{}{
			// Basic mapping - customize based on your Brincr API
			"sku":         safeFetch(wooProduct, "sku", ""),
			"name":        safeFetch(wooProduct, "name", ""),
			"description": safeFetch(wooProduct, "description", ""),
			// Add more mappings as needed
		}
		
		// Handle price conversion
		if price, ok := wooProduct["price"].(string); ok && price != "" {
			brincrProduct["price"] = price
		} else if regularPrice, ok := wooProduct["regular_price"].(string); ok && regularPrice != "" {
			brincrProduct["price"] = regularPrice
		}
		
		// Stock quantity
		if stock, ok := wooProduct["stock_quantity"].(float64); ok {
			brincrProduct["stock_quantity"] = stock
		}
		
		brincrProducts = append(brincrProducts, brincrProduct)
	}
	
	// Print example transformed product
	if len(brincrProducts) > 0 {
		prettyPrint("Transformed Brincr Product Example:", brincrProducts[0])
	}
	
	// Step 7: Push to Brincr (optional - comment out if not ready to push)
	log.Println("Push to Brincr? (This will create actual products in your Brincr account)")
	log.Println("Waiting 10 seconds - press Ctrl+C to cancel")
	time.Sleep(10 * time.Second)
	
	log.Println("Pushing products to Brincr...")
	for i, product := range brincrProducts {
		log.Printf("Pushing product %d/%d: %s", i+1, len(brincrProducts), product["name"])
		
		// Construct the API path
		path := "/api/v1/products"
		
		// Create the product in Brincr
		var result interface{}
		err := brincrClient.Execute("POST", path, product, &result)
		if err != nil {
			log.Printf("Warning: Failed to push product to Brincr: %v", err)
			continue
		}
		
		log.Printf("Successfully pushed product to Brincr")
	}
	
	log.Println("Integration test completed!")
}

// checkRequiredEnvVars ensures all needed environment variables are set
func checkRequiredEnvVars() {
	required := []string{
		// WooCommerce
		"WOOCOMMERCE_API_URL",
		"WOOCOMMERCE_CONSUMER_KEY",
		"WOOCOMMERCE_CONSUMER_SECRET",
		// Brincr
		"BRINCR_TENANT_ID",
		"BRINCR_CLIENT_ID",
		"BRINCR_CLIENT_SECRET",
	}
	
	missing := []string{}
	for _, env := range required {
		if os.Getenv(env) == "" {
			missing = append(missing, env)
		}
	}
	
	if len(missing) > 0 {
		log.Println("Please set the following environment variables:")
		for _, env := range missing {
			log.Printf("- %s\n", env)
		}
		log.Fatalln("Missing required environment variables")
	}
}

// prettyPrint formats and prints JSON data
func prettyPrint(label string, data interface{}) {
	fmt.Println("\n" + label)
	fmt.Println("----------------------------------------")
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
	fmt.Println("----------------------------------------\n")
}

// safeFetch safely extracts values from a map with a default fallback
func safeFetch(data map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, ok := data[key]; ok && value != nil {
		return value
	}
	return defaultValue
}
