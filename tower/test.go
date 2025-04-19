package main

import (
    "log"

    "github.com/joho/godotenv"
    "github.com/zepzeper/tower/internal/connectors/client"
)

type result interface{}

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    clientFactory := client.NewFactory();
    brincrClient, err := clientFactory.CreateClient("brincr", true)
    if err != nil {
        log.Fatalf("Failed to create Brincr client: %v", err)
    }

    // Make a test request
    testResult, err := brincrClient.TestRequest()
    if err != nil {
        log.Printf("Test request failed: %v", err)
    }

    log.Printf("Test request success: %v", &testResult);
}
