package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/zepzeper/tower/internal/config"
	"github.com/zepzeper/tower/internal/connectors/woocommerce"
	"github.com/zepzeper/tower/internal/helpers"
)

func main() {
  helpers.LoadEnvFile()
	// Ensure environment variables are set
	checkEnvVars()

	// Create a new WooCommerce client
	client, err := woocommerce.NewClient(false)
	if err != nil {
		log.Fatalf("Failed to create WooCommerce client: %v", err)
	}

	// Test the connection
	log.Println("Testing connection to WooCommerce...")
	testResult, err := client.TestRequest()
	if err != nil {
		log.Fatalf("Connection test failed: %v", err)
	}
	log.Println("Connection test successful!")

	fmt.Println("WooCommerce Test Result:")
	jsonData, err := json.MarshalIndent(testResult, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
	} else {
		fmt.Println(string(jsonData))
	}

	log.Println("Test completed successfully")
}

// checkEnvVars ensures the required environment variables are set
func checkEnvVars() {
	required := []string{
		"WOOCOMMERCE_API_URL",
		"WOOCOMMERCE_CONSUMER_KEY",
		"WOOCOMMERCE_CONSUMER_SECRET",
	}

	missing := []string{}
	for _, env := range required {
		if config.GetEnv(env, "") == "" {
			missing = append(missing, env)
		}
	}

	if len(missing) > 0 {
		log.Println("The following environment variables must be set:")
		for _, env := range missing {
			log.Printf("- %s\n", env)
		}
		log.Fatalln("Please set the required environment variables and try again")
	}
}

