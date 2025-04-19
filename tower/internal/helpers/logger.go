package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Entry defines the structure of a JSON log entry
type Entry struct {
	URI      string      `json:"uri"`
	Response interface{} `json:"response"`
}

// prettyPrint formats and prints JSON data for readability
func (e *Entry) PrettyPrint(label string, data interface{}) {
	fmt.Println("\n" + label)
	fmt.Println("----------------------------------------")
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
	fmt.Println("----------------------------------------")
}

func LoadEnvFile() {
	// Try to find .env in current directory
	err := godotenv.Load()
	if err != nil {
		// If not found, try to find it in the project root
		cwd, err := os.Getwd()
		if err == nil {
			// Try various possible locations for .env
			possiblePaths := []string{
				filepath.Join(cwd, ".env"),
				filepath.Join(cwd, "../.env"),
				filepath.Join(cwd, "../../.env"),
			}

			for _, path := range possiblePaths {
				if _, err := os.Stat(path); err == nil {
					godotenv.Load(path)
					log.Printf("Loaded environment from %s", path)
					break
				}
			}
		}
	}
}
